# Croney

## Rationale
Occasionally you want to run scheduled tasks that need to execute on a schedule. Doing this on an unknown set of machines can be difficult, requiring a particular machine to be nominated as the scheduler/worker etc.

Croney solves the problem generally: allowing a command to be executed on a regular cron schedule but self-coordinating with other machines to ensure only a single task is run.
