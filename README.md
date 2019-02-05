<h2>get-pull-requests</h2>
This is a CLI tool to list pull requests in your organizations. It will allow you to choose an organization you are interested in or look at all your organizations. The only command available is pulls. The --organization flag allows you to specify which organization you are interested in.


```
NAME:
   Github Pull Request Lookup CLI - Let's you look up for pull requests within your organizations

USAGE:
   get-pull-requests [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     pulls              Looks up the pull requests in your organizations
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
   ```
<h3>Installation</h3>   
Go to releases and download the binary. Open the terminal and chmod 755 the binary to make it executable. You can now use the script. For example, if you were to run it from your downloads, you would run:

```
/Users/$USER/Downloads/get-pull-requests pulls
   ```

or

```
/Users/$USER/Downloads/get-pull-requests pulls --organization dummyorganization123
   ```
