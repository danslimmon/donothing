# donothing

`donothing` is a Go framework for [do-nothing
scripting](https://blog.danslimmon.com/2019/07/15/do-nothing-scripting-the-key-to-gradual-automation/).
Do-nothing scripting is an approach to task automation that lets your processes evolve with minimal
effort from documented manual processes to fully automated.

A `donothing` script walks through a **procedure**, which comprises a sequence of **steps**. As an
example, here's a simple procedure for restoring a database backup:

```go
package main

import (
  "github.com/danslimmon/donothing"
)

func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)

  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
  })
  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("loadBackupData")
    pcd.Short("Load the backup data into the database")
  })
  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("testRestoredData")
    pcd.Short("Check the restored data for consistency")
  })
}
```

When this code is run, the user will be prompted to follow the instructions step by step:

```markdown
# Restore a database backup

## Retrieve the backup file

Press `Enter` when done:

## Load the backup data into the database

Press `Enter` when done:

## Check the restored data for consistency

Press `Enter` when done:

Done!
```

Details for each step may be added in template files. For example, if we put the following markdown
into the file `templates/retrieveBackupFile.md`,

```markdown
- Log in to the AWS console
- Navigate to the `db_backups` S3 bucket
- Download the `.sql` file you wish to restore
```

then, when we run our code, the output will look like this:

```markdown
# Restore a database backup

## Retrieve the backup file

- Log in to the AWS console
- Navigate to the `db_backups` S3 bucket
- Download the `.sql` file you wish to restore

Press `Enter` when done:
```

The main idea behind `donothing` is that, when you're ready to automate a step instead of performing
it manually, you can just add a `Run()` function to the step. Continuing with the example above,
suppose we write a `retrieveBackupFile` function that downloads the latest database backup to our
working directory. We can then automate the first step of our procedure:

```go
// retrieveBackupFile downloads the latest backup file from the S3 bucket.
//
// It returns the path to the local file containing the backup data.
func downloadBackupFile() (string, error) {
  // ... use the AWS API to download the latest backup file ...
  return "backup.sql", nil
}

func main() {
  // ...
  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
    pcd.Run(func(pcd donothing.Procedure)) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      pcd.OutputString("backupFilePath", filename)
    }
  })
  // ...
}
```

Now when we run our script, it will automatically download the backup file from S3, allowing us to
move on to the second step immediately:

```markdown
# Restore a database backup

## Retrieve the backup file

    backupFilePath: ./backup.sql

## Load the backup data into the database

Press `Enter` when done:
```

This paradigm makes it easy to automate the restore procedure piece by piece.

## Inputs and outputs

Sometimes we need to pass information from one step to another. We can do this with step **inputs**
and **outputs**. Continuing with the database restore example, we can have the first step pass the
name of our backup file to the second step, so that the second step can print it.

First, we add a template for the `loadBackupData` step by putting the following Markdown into
`templates/loadBackupData.md`:

```markdown
Run the following command:

    load_data.sh < {{.Input "backupFilePath"}}
```

With no further changes, our `main` function currently looks like this:

```go
func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)

  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
    pcd.Run(func(pcd donothing.Procedure)) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      pcd.OutputString("backupFilePath", filename)
    }
  })

  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("loadBackupData")
    pcd.Short("Load the backup data into the database")
  })

  // ... further steps
}
```

The output from our script will now be:

```markdown
# Restore a database backup

## Retrieve the backup file

    backupFilePath: ./backup.sql

## Load the backup data into the database

Run the following command:

    load_data.sh < ./backup.sql

Press `Enter` when done:
```

Now the user doesn't have to construct their own command for the `loadBackupData` step: they can
just copy and paste the command they need to run. And when it comes time to automate the
`loadBackupData` step, our new `Run` function can use the `backupFilePath` input as well:

```go
func main() {
  // ...

  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
    pcd.Run(func(pcd donothing.Procedure)) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      pcd.OutputString("backupFilePath", filename)
    }
  })

  pcd.AddStep(func(pcd donothing.Procedure) {
    pcd.Name("loadBackupData")
    pcd.Short("Load the backup data into the database")
    pcd.Run(func(pcd donothing.Procedure)) error {
      backupFilePath, ok := pcd.InputString("backupFilePath")
      if !ok {
        return errors.New("Missing required input `backupFilePath`")
      }
      err := loadBackupData(backupFilePath)
      if err != nil {
        return err
      }
      pcd.Log("Data loaded successfully.")
    }
  })

  // ...
}
```

The output from our script will now look like this:

```markdown
# Restore a database backup

## Retrieve the backup file

    backupFilePath: ./backup.sql

## Load the backup data into the database

Data loaded successfully.

    load_data.sh < ./backup.sql

## Check the restored data for consistency

Press `Enter` when done:
```
