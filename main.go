package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"net/http"
	"runtime"
	"sync"

	"encoding/json"
	"code.cloudfoundry.org/cli/plugin"
	"code.cloudfoundry.org/cli/cf/appfiles"
	"code.cloudfoundry.org/cli/cf/terminal"
	"code.cloudfoundry.org/cli/cf/trace"
	"github.com/parnurzeal/gorequest"
	"github.com/simonleung8/flags"
	"github.com/spf13/viper"
	"github.com/xiwenc/cf-fastpush-controller/lib"
	"github.com/xiwenc/cf-fastpush-controller/utils"
	"io/ioutil"
	"strings"
)

const windowsPathPrefix = `\\?\`

/*
*	This is the struct implementing the interface defined by the core CLI. It can
*	be found at  "github.com/cloudfoundry/cli/plugin/plugin.go"
*
 */
type FastPushPlugin struct {
	ui terminal.UI
}

type VCAPApplication struct {
	VCAP_APPLICATION struct {
		ApplicationID      string `json:"application_id"`
		ApplicationVersion string `json:"application_version"`
	} `json:"VCAP_APPLICATION"`
}

type FileEntry struct {
	Checksum     string
	Modification int64
	Content      []byte
}

var lock = sync.RWMutex{}
var store = map[string]*FileEntry{}

/*
*	This function must be implemented by any plugin because it is part of the
*	plugin interface defined by the core CLI.
*
*	Run(....) is the entry point when the core CLI is invoking a command defined
*	by a plugin. The first parameter, plugin.CliConnection, is a struct that can
*	be used to invoke cli commands. The second paramter, args, is a slice of
*	strings. args[0] will be the name of the command, and will be followed by
*	any additional arguments a cli user typed in.
*
*	Any error handling should be handled with the plugin itself (this means printing
*	user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
*	1 should the plugin exits nonzero.
 */
func (c *FastPushPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		fmt.Println("See you.")
		return
	}
	// Ensure that the user called the command fast-push
	// alias fp is auto mapped
	var dryRun bool
	traceLogger := trace.NewLogger(os.Stdout, true, os.Getenv("CF_TRACE"), "")
	c.ui = terminal.NewUI(os.Stdin, os.Stdout, terminal.NewTeePrinter(os.Stdout), traceLogger)

	cliLogged, err := cliConnection.IsLoggedIn()
	if err != nil {
		c.ui.Failed(err.Error())
	}

	if cliLogged == false {
		c.ui.Failed("Cannot perform fast-push without being logged in to CF")
		os.Exit(1)
	}

	// set flag for dry run
	fc := flags.New()
	fc.NewBoolFlag("dry", "d", "bool dry run flag")

	err = fc.Parse(args[1:]...)
	if err != nil {
		c.ui.Failed(err.Error())
	}
	// check if the user asked for a dry run or not
	if fc.IsSet("dry") {
		dryRun = fc.Bool("dry")
		args = args[:len(args)-1]
	}
	apps := []string{}
	if len(args) <= 1 {
		viper.SetConfigName("manifest")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		err := viper.ReadInConfig()
		if err != nil {
			c.ui.Failed(err.Error())
			os.Exit(1)
		}
		applications, ok := viper.Get("applications").([]interface{})
		if !ok {
			c.ui.Failed("manifest parse error.: applications")
			os.Exit(1)
		}
		for _, application := range applications {
			command := application.(map[interface{}]interface{})["command"]
			if command != nil && strings.Contains(command.(string), "./fp") {
				apps = append(apps, application.(map[interface{}]interface{})["name"].(string))
			}
		}
		if apps == nil {
			c.ui.Failed("app not found include fast-push-controller in manifest.yml")
			os.Exit(0)
		}
	} else {
		apps = args[1:]
	}
	if args[0] == "fast-push" || args[0] == "fp" {
		for _, app := range apps {
			c.ui.Say("Running the fast-push command")
			c.ui.Say("Target app: %s \n", app)
			c.FastPush(cliConnection, app, dryRun)
		}
	} else if args[0] == "fast-push-status" || args[0] == "fps" {
		for _, app := range apps {
			c.FastPushStatus(cliConnection, app)
		}
	} else {
		return
	}

}

func (c *FastPushPlugin) GetAuthToken(cliConnection plugin.CliConnection, appName string) string {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		c.ui.Failed(err.Error())
		os.Exit(1)
	}
	return app.Guid
}

func (c *FastPushPlugin) FastPushStatus(cliConnection plugin.CliConnection, appName string) {
	authToken := c.GetAuthToken(cliConnection, appName)

	apiEndpoint := c.GetApiEndpoint(cliConnection, appName)
	status := lib.Status{}
	request := gorequest.New()
	_, body, err := request.Get(apiEndpoint+"/status").Set("x-auth-token", authToken).End()
	if err != nil {
		c.ui.Failed(fmt.Errorf("request.Get error: %s/status.%s", apiEndpoint, err).Error())
		os.Exit(1)
	}
	json.Unmarshal([]byte(body), &status)
	if status.Health != "" {
		c.ui.Say(appName + ": fast-push-controller is " + status.Health)
	} else {
		c.ui.Say(appName + ": fast-push-controller is not running.")
	}
}

func (c *FastPushPlugin) FastPush(cliConnection plugin.CliConnection, appName string, dryRun bool) {
	// Please check what GetApp returns here
	// https://github.com/cloudfoundry/cli/blob/master/plugin/models/get_app.go

	authToken := c.GetAuthToken(cliConnection, appName)

	apiEndpoint := c.GetApiEndpoint(cliConnection, appName)
	request := gorequest.New()
	response, body, err := request.Get(apiEndpoint+"/files").Set("x-auth-token", authToken).End()
	if err != nil {
		c.ui.Failed(fmt.Errorf("request.Get error: %s/files: %s", apiEndpoint, err).Error())
		os.Exit(1)
	}
	if response.StatusCode != http.StatusOK {
		c.ui.Failed(fmt.Errorf("response code is not 200 ok: %s/files. code=%s. %s", apiEndpoint, response.StatusCode, err).Error())
		os.Exit(1)
	}
	remoteFiles := map[string]*lib.FileEntry{}
	json.Unmarshal([]byte(body), &remoteFiles)

	localFiles := ListFiles()

	filesToUpload := c.ComputeFilesToUpload(localFiles, remoteFiles)
	payload, _ := json.Marshal(filesToUpload)
	if dryRun {
		// NEED TO HANDLE DRY RUN
		c.ui.Warn("warn: No changes will be applied, this is a dry run.")
	} else {

		_, body, err = request.Put(apiEndpoint+"/files").Set("x-auth-token", authToken).Send(string(payload)).End()
		if err != nil {
			c.ui.Failed(fmt.Errorf("request.Get error: %s/files: %s", apiEndpoint, err).Error())
			os.Exit(1)
		}
	}

	status := lib.Status{}
	json.Unmarshal([]byte(body), &status)
	c.ui.Say(status.Health)
}

/*
*	This function must be implemented as part of the plugin interface
*	defined by the core CLI.
*
*	GetMetadata() returns a PluginMetadata struct. The first field, Name,
*	determines the name of the plugin which should generally be without spaces.
*	If there are spaces in the name a user will need to properly quote the name
*	during uninstall otherwise the name will be treated as seperate arguments.
*	The second value is a slice of Command structs. Our slice only contains one
*	Command Struct, but could contain any number of them. The first field Name
*	defines the command `cf basic-plugin-command` once installed into the CLI. The
*	second field, HelpText, is used by the core CLI to display help information
*	to the user in the core commands `cf help`, `cf`, or `cf -h`.
 */
func (c *FastPushPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "cf-fastpush",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 3,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 28,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "fast-push",
				Alias:    "fp",
				HelpText: "fast-push removes the need to deploy your app again for a small change",
				UsageDetails: plugin.Usage{
					Usage: "cf fast-push APP_NAME\n   cf fp APP_NAME",
					Options: map[string]string{
						"dry": "--dry, dry run for fast-push",
					},
				},
			},
			plugin.Command{
				Name:     "fast-push-status",
				Alias:    "fps",
				HelpText: "fast-push-status shows the current state of your application",
				UsageDetails: plugin.Usage{
					Usage: "cf fast-push-status APP_NAME\n   cf fps APP_NAME",
				},
			},
		},
	}
}

/*
* Unlike most Go programs, the `Main()` function will not be used to run all of the
* commands provided in your plugin. Main will be used to initialize the plugin
* process, as well as any dependencies you might require for your
* plugin.
 */
func main() {
	plugin.Start(new(FastPushPlugin))
}

func (c *FastPushPlugin) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}

func (c *FastPushPlugin) GetApiEndpoint(cliConnection plugin.CliConnection, appName string) string {
	results, err := cliConnection.CliCommandWithoutTerminalOutput("app", appName)
	if err != nil {
		c.ui.Failed(err.Error())
	}
	// for debug
	/*
		for _, line := range results {
			fmt.Println(line)
		}
	*/
	app_env, err := cliConnection.CliCommandWithoutTerminalOutput("env", appName)
	if err != nil {
		c.ui.Failed(err.Error())
	}

	protocol := "https"
	rep := regexp.MustCompile("^FP_PROTOCOL: (.+)")
	for _, line := range app_env {
		if rep.MatchString(line) {
			protocol = strings.ToLower(strings.TrimSpace(rep.ReplaceAllString(line, "$1")))
			if protocol != "http" && protocol != "https" {
				c.ui.Failed("invalid FP_PROTOCOL: " + protocol)
				os.Exit(1)
			}
			break
		}
	}

	fp_domain := ""
	rep = regexp.MustCompile("^FP_DOMAIN: (.+)")
	for _, line := range app_env {
		if rep.MatchString(line) {
			fp_domain = strings.ToLower(strings.TrimSpace(rep.ReplaceAllString(line, "$1")))
			break
		}
	}

	fp_endpoint := ""
	rep = regexp.MustCompile("(?i)^(?:adresses )?url[s]?[ ]?:.*[[:space:]](.*\\." + fp_domain + ").*$")
	for _, line := range results {
		if rep.MatchString(line) {
			fp_endpoint = protocol + "://" + strings.TrimSpace(rep.ReplaceAllString(line, "$1")) + "/_fastpush"
			break
		}
	}
	if fp_endpoint == "" {
		c.ui.Failed("Could not find usable route for this app. Make sure at least one route is mapped to this app.")
		os.Exit(1)
	}

	app_instances := ""
	rep = regexp.MustCompile("(?i)^(?:instances|instances |instanzen|instancias|インスタンス|istanze|인스턴스|instâncias|实例|實例) ?: (.+)$")
	for _, line := range results {
		if rep.MatchString(line) {
			app_instances = strings.TrimSpace(rep.ReplaceAllString(line, "$1"))
			break
		}
	}
	if app_instances != "1/1" {
		c.ui.Failed("The number of instances of app must be 1. Current: " + app_instances)
		os.Exit(1)
	}

	fmt.Println("Endpoint: " + fp_endpoint)
	return fp_endpoint
}

func (c *FastPushPlugin) ComputeFilesToUpload(local map[string]*FileEntry, remote map[string]*lib.FileEntry) map[string]*FileEntry {
	filesToUpload := map[string]*FileEntry{}
	for path, f := range local {
		if remote[path] == nil {
			c.ui.Say("[NEW] " + path)
			f.Content, _ = ioutil.ReadFile(path)
			filesToUpload[path] = f
		} else if remote[path].Checksum != f.Checksum {
			c.ui.Say("[MOD] " + path)
			f.Content, _ = ioutil.ReadFile(path)
			filesToUpload[path] = f
		}
	}
	return filesToUpload
}

func ListFiles() map[string]*FileEntry {
	dir := "./"
	cfIgnore := loadIgnoreFile(dir)
	log.Println("Listing files for: " + dir)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileRelativePath, _ := filepath.Rel(dir, path)
		fileRelativeUnixPath := filepath.ToSlash(fileRelativePath)
		if err != nil && runtime.GOOS == "windows" {
			f, err = os.Lstat(windowsPathPrefix + path)
			if err != nil {
				return err
			}
			path = windowsPathPrefix + path
		}

		if f.IsDir() {
			return nil
		}
		if cfIgnore.FileShouldBeIgnored(fileRelativeUnixPath) {
			if err == nil && f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if store[path] != nil && store[path].Modification == f.ModTime().Unix() {
			// cache hit
			return nil
		}

		fileEntry := FileEntry{}
		checksum, _ := utils.ChecksumsForFile(path)
		fileEntry.Checksum = checksum.SHA256
		fileEntry.Modification = f.ModTime().Unix()
		lock.RLock()
		store[path] = &fileEntry
		lock.RUnlock()
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return store
}

func loadIgnoreFile(dir string) CfIgnore {
	fileContents, err := ioutil.ReadFile(filepath.Join(dir, ".cfignore"))
	if err != nil {
		return appfiles.NewCfIgnore("")
	}

	return appfiles.NewCfIgnore(string(fileContents))
}

type CfIgnore interface {
	FileShouldBeIgnored(path string) bool
}
