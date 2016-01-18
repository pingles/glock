# Croney

## Rationale
Occasionally you want to run tasks that need to execute on a schedule. Like cron ;-)

However, doing this on an unknown set of machines is a little more difficult if you want only one instance of the task to execute. Normally one machine is the special one: tasked as the scheduler, worker etc.

Croney solves the problem generally: allowing a command to be executed on a regular cron schedule but self-coordinating with other machines to ensure only a single task is run.
