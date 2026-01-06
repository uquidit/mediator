# uQuid-IT Mediator

# Context

uQuidIT.co provides this mediator script system to customize your workflow in Tufin TOS Aurora.

## General purpose

This system allows you to execute any external script to process your Securechange tickets. It is a convenient solution to use your legacy scripts developed for TOS Classic in an Aurora environment. 

These scripts will be executed on the main node instead of the Securechange pod, providing extra flexibility.

We assume you are familiar with Tufin Securechange and know how to configure it.
## Architecture overview

The main difficulty to overcome is that in the TOS Classic architecture, scripts are executed on the machine where SecureChange runs while in the TOS Aurora architecture, those scripts are supposed to run on a Kubernetes pod, which is very limited.

The solution is to deploy only a very small, self-contained executable on the pod: the `mediator-client`. This client will communicate with a `mediator-server`, which is hosted on a server that can execute the legacy scripts on behalf of the client.

The following diagram illustrates the interactions of the different elements involved.

![Architecture overview](mediator_client_and_server_.png)

# Build

## Get the code

You need to download the source code to build the executable file. 

To do so, clone our Github repository:
`$ git clone https://github.com/uquidit/mediator.git`

## Get Golang

You need Golang installed so you can build the binaries. Follow [this procedure](https://go.dev/doc/install) from the [go.dev](go.dev) web site.

## Build

3 binaries should be compiled:
* `mediator-client`
* `mediator-server`
* `mediator-cli`

You can compile them by going to their respective directories and run the `go build .` command. The binary will be created in current directory. Use the `-o` flag to specify a destination directory.

* `mediator-client`  => `mediator/cmd/mediator-client`
* `mediator-server`  => `mediator/cmd/mediator-server`
* `mediator-cli`     => `mediator/cmd/mediator-cli`

### Set encryption keys

If you use mediator in production environment, you must change the following encryption keys:

Encryption algorithm uses 3 variables in `mediatorscript` package:
* `salt`
* `pepper`
* `secretKey`

And the One-Time password algorithm uses 2 in `totp` package:
* `secretMS1`
* `secretMS2`

They can all be set at build time using the `--ldflags` build flag with the `-X` command: 
`--ldflags="-X '<project name>/<package name>.<variable name>=<string value>'"`

Example:
```sh
$ go build --ldflags="\
-X 'mediator/mediatorscript.salt=Hello' \
-X 'mediator/mediatorscript.pepper=World' \
-X 'mediator/mediatorscript.secretKey=ThisIsMySecret' \
-X 'mediator/totp.secretMS1=Lorem' \
-X 'mediator/totp.secretMS2=Ipsum' \
" .
```

Note: you can use escaped double-quotes to use spaces if required.

### The easy way

There is a make target that takes care of everything for you. Simply run:
```sh
$ make package
```
And you are good to deploy the `package/mediator.tar.gz` artifact!

# Installation & Configuration

Mediator executables must be installed on your server.  

## Mediator-server

### Files

The `mediator-server` is composed of 3 files:

* The main executable file: `mediator-server`;
* The distribution configuration file: `mediator-server_dist.yml`. which will be used as a starting point for our configuration. It provides the comprehensive list of available configuration pieces along with their documentation;
* The Command-Line Interface (CLI) executable file: `mediator-cli` used to manage the server.

### Configuration

We provide a distribution configuration file named `mediator-server_dist.yml` as a starting point to build your own configuration.

You will find a comprehensive list of documented configuration entries.

We recommend that you keep an untouched copy of the original version of the distribution configuration file.

Follow these steps to configure the `mediator-server`:
1. Copy the distribution configuration file: `$ cp mediator-server_dist.yml mediator-server.yml`

2. Edit `mediator-server.yml` with your favorite text editor and set the required values according to your infrastructure and save the file

### Installation

1. Create the following destination folders on your server and make sure they are read and writeable for the user that will run `mediator-server` (likely to be the `tufin-admin` user)
   * `/opt/mediator/conf`
   * `/opt/mediator/bin`
   * `/opt/mediator/data`
2. Copy executable files to `/opt/mediator/bin` folder and set executable flag
3. Copy the configuration file to `/opt/mediator/conf` folder

### Start `mediator-server`

If you want the `mediator-server` to bind on a privileged port, it's necessary to grant the related special capability to it's binary file. This must be done after each installation as well as update. Run the following command:

`$ sudo setcap 'cap_net_bind_service=+ep' /opt/mediator/bin/mediator-server`

`mediator-server` expects the full path to its configuration file as the only argument.

If you’re testing the mediator, we recommend that you start the server in a “screen”:

```
$ screen -S mediator
$ /opt/mediator/bin/mediator-server /opt/mediator/conf/mediator-server.yml
```

Hit “Ctrl+a” then “d” to detach the screen without killing it. Use the following command to reattach:

`$ screen -r mediator`

If you’re installing the mediator on a production server, we recommend that you start it using systemd or other demon management software available on your system.

If you’re using systemd – which is likely to be – you can copy the provided file `cmd/mediator-server/mediator-server.service` into the `/etc/systemd/system` directory on the server:

```
$ sudo cp /tmp/mediator/mediator-server.service /etc/systemd/system
$ sudo systemctl daemon-reload
$ sudo systemctl enable mediator-server.service
```

You can check the service is enabled:
```
$ sudo systemctl is-enabled mediator-server.service
enabled
```

You can manage your service with the following commands:

* `sudo systemctl start mediator-server.service`
* `sudo systemctl stop mediator-server.service`
* `sudo systemctl status mediator-server.service`
* `sudo systemctl restart mediator-server.service`

### Command-line interface

The command-line interface is a convenient utility to easily interact with the mediator server. In order to use it, you need to provide the URL the mediator server listens to as the command line utility has no configuration file.

We suggest that you install it in the same folder as the server but it is not a requirement. In the following section, we will assume this is the case.

For instance, the following command will return the list of available commands:

```
$ /opt/mediator/bin/mediator-cli --url=http://xxx/v1/otp --help

This CLI provides commands to manage scripts used by Mediator back-end. It includes:
* List registered scripts
* Register new scripts
* Un-register useless scripts
* Refresh script checksum
* Test scripts

Usage:
  mediator [command]

Available Commands:
  completion       Generate the autocompletion script for the specified shell
  help             Help about any command
  script           List available scripts for mediator. Available alias:'scripts'
  securechange-api Show and manage SecurechangeAPI configuration
  settings         Update Mediator settings

Flags:
  -h, --help            help for mediator
      --sslskipverify   Skip SSL certificate verification (insecure)
  -u, --url string      Back-end URL (required)

Use "mediator [command] --help" for more information about a command.
```

The `--help` flag will always give you tips about how to use the CLI. 

A help message is available for all commands and subcommands. 

An alias may come handy when using the CLI so you don’t have to type the long server URL every time. There are different ways to do so, here is one of them using the alias command:

```
$ alias mediator="/opt/mediator/bin/mediator-cli --url=http://xxx/v1/otp"
```

In the following sections, we will use this alias in the commands we provide.

### Final configuration: script registration

`mediator-client` will request `mediator-server` to run some scripts when a ticket is processed under Securechange. Those scripts must explicitely be registered to the `mediator-server`. We will use `mediator-cli` for that purpose. 

4 different types of scripts can be used:

* **Trigger Scripts:**
  - `mediator-client` can be called when a workflow action is triggered. For instance, when a ticket advances or is submitted.
  - It will request the execution of any Trigger Script attached to the workflow step the ticket is in.
  - Mediator supports use of multiple Trigger Scripts so you can use different scripts to interact with your tickets at different steps of different workflows.
  - You can only attach one Trigger Script to a workflow step.
  - They run asynchronously
  - You can manage them using the `mediator scripts trigger` command and subcommands

* **“interactive” scripts**
  - They are 3 of them:
    - Scripted Condition scripts
    - Scripted Task scripts
    - Pre-Assignment scripts
  - `mediator-client` can be called when a ticket is in a step where a _scripted condition_, a _scripted task_ or _pre-assignment_ is required
  - It will  request the execution of any registered script of the corresponding type.
  - Only one script of each type can be registered on the server.
  - They run synchronously
  - You can manage them using the corresponding command and subcommands:
    - `mediator scripts condition [...]`
    - `mediator scripts task [...]`
    - `mediator scripts assignment [...]`

The commands `mediator scripts [trigger|condition|task|assignment]` all offer 4 subcommands:
* `refresh`
* `register`
* `test`
* `unregister`

All these subcommands work the same way regardless of the type of script.

All the scripts need to be registered before they can be executed by the server. You can use the following command for that purpose:

```
mediator scripts [trigger|condition|task|assignment] register
```

The command requires the script path as a unique argument. 

For instance, to register the “run.sh” script as a Trigger Script:

```
$ mediator scripts trigger register /path/to/script/run.sh
```

A new Trigger script named "run.sh" has been registered and linked to file `/path/to/script/run.sh`

All the scripts are registered under a name. This name must be unique in the script database. By default, the mediator will use the script file name as a name. However, you can provide a custom name using the `--name` flag.

Example:
```
$ mediator scripts trigger register /path/to/scripts/my-script.py --name MyScript
```

A new Trigger script named "MyScript" has been registered and linked to file `/path/to/scripts/my-script.py`

Script registration names are useful for trigger scripts. You will need the script name when you attach it to a workflow step in `mediator-client` configuration.

Check the scripts have been properly registered using the following command:

```
$ mediator scripts trigger 
Nb of Trigger script: 2
- run.sh: /path/to/scripts/run.sh
- MyScript: /path/to/scripts/my-script.py
```

You can unregister a script with the `unregister` command:
```
$ mediator scripts trigger unregister MyScript
Script 'MyScript' has been unregistered.

$ mediator scripts trigger 
Nb of Trigger script: 1
- run.sh: /path/to/scripts/run.sh
```

You can also unregister all the scripts of a given type by using the `unregister` subcommand with no argument

For security reasons, the server will not run a script that has been modified after it has been registered. If you do need to update a registered script, please run the following command after the script has been modified:

```
$ mediator script trigger refresh run.sh
Script ‘run.sh’ has been refreshed
```

All the subcommands described above are available for the other script types. Just change the trigger subcommand by the corresponding command name.

Top-level subcommands are also available. They will operate on all scripts, regardless of their type. Use the `--help` flag for more information.












## Mediator-client installation


### Helper scripts

`mediator-client` comes with 13 helper scripts. They embed extra configuration and make it easy to use and configure `mediator-client`. In facts, `mediator-client` executable will never be run directly.

Make sure they are available on the main node:

* `mediator-client-advance.sh`: executes `mediator-client` when the trigger `ADVANCE` is fired.
* `mediator-client-automatic-step-failed.sh`: executes `mediator-client` when the trigger `AUTOMATIC STEP FAILED` is fired.
* `mediator-client-cancel.sh`: executes `mediator-client` when the trigger `CANCEL` is fired.
* `mediator-client-close.sh`: executes `mediator-client` when the trigger `CLOSE` is fired.
* `mediator-client-create.sh`: executes `mediator-client` when the trigger `CREATE` is fired.
* `mediator-client-redo.sh`: executes `mediator-client` when the trigger `REDO` is fired.
* `mediator-client-reject.sh`: executes `mediator-client` when the trigger `REJECT` is fired.
* `mediator-client-reopen.sh`: executes `mediator-client` when the trigger `REOPEN` is fired.
* `mediator-client-resolve.sh`: executes `mediator-client` when the trigger `RESOLVE` is fired.
* `mediator-client-resubmit.sh`: executes `mediator-client` when the trigger `RESUBMIT` is fired.
* `mediator-client-pre-assignment.sh`: To be used as "Pre-Assignment" script
* `mediator-client-scripted-condition.sh`: To be used as "Scripted condition" script
* `mediator-client-scripted-task.sh`: To be used as "Scripted task" script

### Upload to Securechange pod

15 files need to be uploaded to Securechange:
* `mediator-client` executable file and its configuration file:
  - `mediator-client`
  - `mediator-client.yml`: see next chapter to know how to create this file.
* 13 `mediator-client` helpers listed above

*NB: `mediator-client` also needs a settings file but it will be automatically uploaded to Securechange. You do not need to deal with it at this stage.*

Upload them all to Securechange using the following command for each of the 15 previous files as argument:

```
$ sudo tos scripts sc push mediator-client
[Mar 20 23:59:57]  INFO Pushing from "mediator-client" to "."
[Mar 20 23:59:57]  INFO Done pushing files/folders
```




## Mediator-client: Configuration & Settings

The configuration of `mediator-client` includes information about how it should connect to `mediator-server` and log its activity. 
`mediator-client` also requires a setting file which contains the list of scripts that need to be executed when a ticket reaches a particular step of its workflow. 

Configuration and settings are split into 2 different files:
* A YAML file `mediator-client.yml` contains back-end connection configuration
* A JSON file `mediator-client.json` contains information about script execution. *It should not be manually edited*.

### Configuration file: how mediator-client connects to server

A distribution `mediator-client_dist.yml` YAML file is provided for your convenience:

```yaml
configuration:
  backend_url: https://DOMAIN/v1/otp/mediatorscript
  ssl_skip_verify: false
  log:
    file: /var/log/mediator-client.log
    level: info
```

Steps:
1. Copy `mediator-client_dist.yml` into `mediator-client.yml`
2. Edit the newly created `mediator-client.yml` file and enter your back-end server information and logging preferences.
3. Save the file. *Do NOT change the file name!* 
4. Upload the file to the Securechange pod in the same directory as the `mediator-client` executable.

### Settings file: which script mediator-client should trigger

`mediator-cli` will assist you editing and uploading the settings file to Securechange.

Prior to going into the following steps, make sure that all the required workflows are properly activated in Securechange.

The `settings` command will be used to generate and upload the settings file:


```
$ mediator settings --help
Interactively update Mediator client settings.
	
These settings tells Mediator client which script should be run for a given ticket and trigger.
This is an interactive command: required information will be prompted to you.

Usage:
  mediator settings [flags]

Flags:
  -h, --help              help for settings
  -H, --host string       SecureChange host. Will be prompted if not provided.
  -P, --password string   SecureChange password. Will be prompted if not provided.
  -s, --settings string   Path to local Mediator client settings file.
  -U, --username string   SecureChange user name. Will be prompted if not provided.

Global Flags:
      --sslskipverify   Skip SSL certificate verification (insecure)
  -u, --url string      Back-end URL (required)

```

Run the command. It will download any existing settings file from Securechange and assist you while editing the file.
If you do not provide connection information via `--host`, `--username` and `--password` flag, they will be prompted at runtime:

```
$ mediator settings 
Get fresh list of workflows from SC...
SecureChange Username : securechange_admin
SecureChange Password : 
                        ---
SecureChange Host : 
                    ---
OK !
Getting settings from server...    OK !
Get list of registered trigger scripts from backend...    OK !

******************************************
Choose a workflow (! indicates workflows without settings):
 - 1: Rule Recertification !
 - 2: Opening Firewall Request !
 - 3: Simple Workflow !
 - 4: Generic / Provisioning !
 - 5: Exit
 - 6: Save and exit
 : 1
```
*NB : The workflows in previous list is purely fictional. When you run the command, `mediator-cli` will display the list of all active workflows in your instance.*


Select a workflow. `mediator-cli` will ask you if you want to add or edit a rule if previous settings exist for this workflow. In our example, no settings were found for selected workflow
```
Do you want to edit a rule or add a new one
 - 1 : New rule
 - 2*: Exit
 : 
```

**Definition**: a `rule` is composed by a `script`, a `trigger` and sometimes and workflow `step`. It will tell `mediator-client` to run the script whenever the selected  trigger is fired by Securechange.
Some triggers can be fired in several workflow steps. If you use these in your rule, you will be ask to select a workflow step.

Rule creation for a "simple" trigger will look like this:
```
Do you want to edit a rule or add a new one
 - 1 : New rule
 - 2*: Exit
 : 1
Choose a trigger
 - 1: Create
 - 2: Close
 - 3: Cancel
 - 4: Reject
 - 5: Resubmit
 - 6: Resolve
 - 7: Advance
 - 8: Redo
 - 9: Reopen
 - 10: Automatic step failed
 : 1
Choose a script from list of registered trigger scripts
 - 1: script1
 - 2: script2
 - 3: script3
 : 3

Do you want to edit a rule or add a new one
 - 1 : Trigger Create fires script script3
 - 2 : New rule
 - 3*: Exit
 : 
```

Rule creation for a trigger that requires a step will look like this:
```
Do you want to edit a rule or add a new one
 - 1 : Trigger Create fires script coucou
 - 2 : New rule
 - 3*: Exit
 : 2
Choose a trigger
 - 1: Create
 - 2: Close
 - 3: Cancel
 - 4: Reject
 - 5: Resubmit
 - 6: Resolve
 - 7: Advance
 - 8: Redo
 - 9: Reopen
 - 10: Automatic step failed
 : 7
Choose a step:
 - 1: Security review recertification
 - 2: Risk policy update process
 : 1
Choose a script from list of registered trigger scripts
 - 1: script1
 - 2: script2
 - 3: script3
 : 2

Do you want to edit a rule or add a new one
 - 1 : Trigger Create fires script script3
 - 2 : Trigger Advance on step Security review recertification fires script script2
 - 3 : New rule
 - 4*: Exit
 : 
```

Add all the required rules for all the workflows you need to add settings for. 

When you're done, select `Save and exit` from workflow list:
```
******************************************
Choose a workflow (! indicates workflows without settings):
 - 1: Rule Recertification !
 - 2: Opening Firewall Request !
 - 3: Simple Workflow !
 - 4: Generic / Provisioning !
 - 5: Exit
 - 6: Save and exit
 : 6
 Exiting...
Sending settings to backend for upload to Securechange...    OK !
```

### Securechange API configuration

`mediator-client` needs to be registered in SecurechangeAPI. This task can be done via the Securechange GUI but can be particularly cumbersome in our situation.

It can also be done by `mediator-cli` through its `securechange-api` command:

```
$ mediator securechange-api --help
If no subcommand is provided, show SecurechangeAPI configuration.

Usage:
  mediator securechange-api [flags]
  mediator securechange-api [command]

Aliases:
  securechange-api, scapi, sc, api

Available Commands:
  create      Create a new SecurechangeAPI trigger configuration
  delete      Delete SecurechangeAPI configuration
  show        Show SecurechangeAPI configuration

Flags:
  -h, --help               help for securechange-api
  -H, --host string        SecureChange host. Will be prompted if not provided.
  -P, --password string    SecureChange password. Will be prompted if not provided.
  -U, --username string    SecureChange user name. Will be prompted if not provided.
  -w, --workflow strings   Comma separated list of workflow names. Only show SecurchangeAPI configuration that includes a workflow in the provided list. Can also be used multiple times.
```

#### Show current configuration

When no subcommand is provided or when the `show` subcommand is used, the Securechange API configuration will be shown.

You can see which `mediator-client` helper will be run when a trigger is fired. (there is one helper per trigger). This may be a long list.

The following example show a fictional piece of configuration where the `mediator-client` helper `mediator-client-advance.sh` is registered for trigger `ADVANCE` for workflow `Rule Recertification` :
```
*** #94 Rule Recertification ADVANCE
  - Execute: Rule Recertification "/opt/tufin/data/securechange/scripts/mediator-client-advance.sh"
  - Trigger groups:
    - Name: trigger Advance
    - Workflow: Rule Recertification
    - Triggers: [ADVANCE]
```

#### Configure a workflow in Securechange API

The easiest and most robust way to configure Securechange is via the fully automatic version of `securechange-api` command of `mediator-cli`.

We suggest that you use it for all your workflows:
```
$ mediator securechange-api create -w "Rule Recertification" --all-triggers
SecureChange Host : <SC host or IP address>
SecureChange Username : securechange_admin
SecureChange Password : 
                        ---
Get fresh list of workflows from SC...
OK !
SecurechangeAPI trigger was created for trigger(s) [Create Close Cancel Reject Resubmit Resolve Advance Redo Reopen Automatic step failed].
```

Just replace "Rule Recertification" by the name of your workflow.

#### Manually configure a workflow in Securechange API

You can also use the previous command without a workflow name to select your workflow from a list.

If you do not provide the `--all-triggers`, you can select the trigger you want.

*We do not recommend this usage. Do it only if you know what you're doing*

