### This table describes some of the common signals:

| Signal Number | Signal Name | Default Action | Description |
|---|---|---|---|
| 1  | SIGHUP  | Terminate | Hang up controlling terminal or process. Sometimes used to reread configuration files. |
| 2  | SIGINT  | Terminate | Interrupt from keyboard (`Ctrl + C`). |
| 3  | SIGQUIT | Dump | Quit from keyboard (`Ctrl + \`). |
| 9  | SIGKILL | Terminate | Forced termination. |
| 15 | SIGTERM | Terminate | Graceful termination. |
| 17 | SIGCHLD | Ignore | Child process exited. |
| 18 | SIGCONT | Continue | Resume process execution. |
| 19 | SIGSTOP | Stop | Stop process execution (`Ctrl + Z`). |

### Process management using Go

Packages 

- os 
- os/exec
- syscall

 provide a lot of useful functionality for interacting with an OS from a Go application.

 - os/exec allows running external shell commands.
 - syscall provides an interface to the low-level OS primitives and allows executing system calls.