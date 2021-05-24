# donothing

`donothing` is a Go framework for [do-nothing
scripting](https://blog.danslimmon.com/2019/07/15/do-nothing-scripting-the-key-to-gradual-automation/).
Do-nothing scripting is an approach to writing procedures. It allows you to start with a documented
manual process and gradually make it better by automating a step at a time. Do-nothing scripting
aims to minimize the [activation energy](https://en.wikipedia.org/wiki/Activation_energy) for
automating steps of a manual procedure.

A `donothing` script walks through a **procedure**, which comprises a sequence of **steps**. As an
example, here's a do-nothing procedure for restoring a database backup:

```go
package main

import (
  "github.com/danslimmon/donothing"
)

func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)
  pcd.Long(`
    This is our procedure for restoring a database backup. To familiarize yourself with our database
    setup, see [the database docs](https://example.com/docs/database.html).
  `)

  pcd.AddStep(func(step *donothing.Step) {
    step.Name("retrieveBackupFile")
    step.Short("Retrieve the backup file")
    step.Long(`
      Log in to the [storage control panel](https://example.com/storage/) and locate the latest file
      of the form "backup_YYYYMMDD.sql". Download that file to your workstation.
    `)
  })

  pcd.AddStep(func(step *donothing.Step) {
    step.Name("loadBackupData")
    step.Short("Load the backup data into the database")
    step.Long(`
      Run this command to load the backup data into the database:

          psql < backup_YYYYMMDD.sql
    `)
  })

  pcd.AddStep(func(step *donothing.Step) {
    step.Name("testRestoredData")
    step.Short("Check the restored data for consistency")
    step.Long(`
      Log in to the database and make sure there are recent records in the events table.
    `)
  })

  pcd.Execute()
}
```

When this code is run, the user will be prompted to follow the instructions step by step:

```
# Restore a database backup

This is our procedure for restoring a database backup. To familiarize yourself with our database
setup, see [the database docs](https://example.com/docs/database.html).

[Enter] to begin:

## Retrieve the backup file

Log in to the [storage control panel](https://example.com/storage/) and locate the latest file
of the form "backup_YYYYMMDD.sql". Download that file to your workstation.

[Enter] when done:

## Load the backup data into the database

Run this command to load the backup data into the database:

    psql < backup_YYYYMMDD.sql

[Enter] when done:

## Check the restored data for consistency

Log in to the database and make sure there are recent records in the events table.

[Enter] when done:

Done!
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
  return filename, nil
}

func main() {
  // ...
  pcd.AddStep(func(step *donothing.Step) {
    step.Name("retrieveBackupFile")
    step.Short("Retrieve the backup file")

    step.Run(func(facts *donothing.Facts) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      facts.SetString("backupFilePath", filename)
      return nil
    })
  })
  // ...
}
```

Note that the `retrieveBackupFile` step's `Long()` call has been removed, since it's obsoleted by
the automation provided by `Run()`. A step _can_ have both `Long()` and `Run()`, but in this case
it doesn't need to.

Now when we run our script, it will automatically download the backup file from S3, allowing us to
move on to the second step immediately:

```markdown
# Restore a database backup

This is our procedure for restoring a database backup. To familiarize yourself with our database
setup, see [the database docs](https://example.com/docs/database.html).

[Enter] to begin:

## Retrieve the backup file

Executing step `retrieveBackupFile` automatically.

**Outputs**:
  - `backupFilePath`: ./backup_20200226.sql

## Load the backup data into the database

Run this command to load the backup data into the database:

    psql < backup_YYYYMMDD.sql

[Enter] when done:
```

This paradigm makes it easy to automate the procedure piece by piece.

## Facts

Sometimes we need to pass information from one step to another. We can do this with **facts**. A
fact is a uniquely named, typed value. Over the course of a procedure, `donothing` collects the
facts produced by steps, and passes those facts along to subsequent steps.

Continuing with the database restore example, we can have the first step pass the name of our backup
file to the second step, so that the second step can print it.

First, we modify the long description of the `loadBackupData` step so that it contains new,
templated instructions.

```go
  pcd.AddStep(func(step *donothing.Step) {
    step.Name("loadBackupData")
    step.Short("Load the backup data into the database")
    step.Long(`
      Run this command to load the backup data into the database:

          psql < {{.Facts.GetString "backupFilePath"}}
    `)
  })
```

With no further changes, our `main` function now looks like this:

```go
func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)
  pcd.Long(`
    This is our procedure for restoring a database backup. To familiarize yourself with our database
    setup, see [the database docs](https://example.com/docs/database.html).
  `)

  pcd.AddStep(func(step *donothing.Step) {
    step.Name("retrieveBackupFile")
    step.Short("Retrieve the backup file")

    step.Run(func(facts *donothing.Facts) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      facts.SetString("backupFilePath", filename)
      return nil
    })
  })

  pcd.AddStep(func(pcd donothing.Procedure) {
    step.Name("loadBackupData")
    step.Short("Load the backup data into the database")
    step.Long(`
      Run this command to load the backup data into the database:

          psql < {{.Facts.GetString "backupFilePath"}}
    `)
  })

  // ... further steps
}
```

The output from our script will now be:

```markdown
# Restore a database backup

This is our procedure for restoring a database backup. To familiarize yourself with our database
setup, see [the database docs](https://example.com/docs/database.html).

[Enter] to begin:

## Retrieve the backup file

Executing step `retrieveBackupFile` automatically.

**Outputs**:
  - `backupFilePath`: ./backup_20200226.sql

## Load the backup data into the database

Run this command to load the backup data into the database:

    psql < ./backup_20200226.sql

[Enter] when done:
...
```

Now the user doesn't have to construct their own command for the `loadBackupData` step: they can
just copy and paste the command they need to run. And when it comes time to automate the
`loadBackupData` step as well, our new `Run` function can use the `backupFilePath` fact by
retrieving it:

```go
func main() {
  // ...
  pcd.AddStep(func(step *donothing.Step) {
    step.Name("loadBackupData")
    step.Short("Load the backup data into the database")
    step.Run(func(facts *donothing.Facts) error {
      // Second return value, discarded here, will be true iff the fact exists.
      backupFilePath, _ := facts.GetString("backupFilePath")
      err := loadBackupData(backupFilePath)
      if err != nil {
        return err
      }
      fmt.Println("Data loaded successfully.")
      return nil
    }
  })
  // ...
}
```

The output from our script will now look like this:

```markdown
# Restore a database backup

This is our procedure for restoring a database backup. To familiarize yourself with our database
setup, see [the database docs](https://example.com/docs/database.html).

[Enter] to begin:

## Retrieve the backup file

Executing step `retrieveBackupFile` automatically.

**Outputs**:
  - `backupFilePath`: ./backup_20200226.sql

## Load the backup data into the database

Executing step `loadData` automatically.

Data loaded successfully.

## Check the restored data for consistency

Log in to the database and make sure there are recent records in the events table.

[Enter] when done:
```

## Generating procedure documentation

`donothing` can print Markdown documentation for a procedure. Going back to our original,
non-automated database restore example, let's add a `--print` flag to our script:

```go
package main

import (
  "os"
  "github.com/danslimmon/donothing"
)

func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)
  pcd.Long(`
    This is our procedure for restoring a database backup. To familiarize yourself with our database
    setup, see [the database docs](https://example.com/docs/database.html).
  `)

  pcd.AddStep(func(step *donothing.Step) {
    step.Name("retrieveBackupFile")
    step.Short("Retrieve the backup file")
  })
  step.AddStep(func(step *donothing.Step) {
    step.Name("loadBackupData")
    step.Short("Load the backup data into the database")
  })
  pcd.AddStep(func(step *donothing.Step) {
    step.Name("testRestoredData")
    step.Short("Check the restored data for consistency")
  })

  if len(os.Args) > 0 && os.Args[1] == "--print" {
    pcd.Render()
  } else {
    pcd.Execute()
  }
}
```

When we invoke our script with the `--print` flag, it will print out our whole procedure as
Markdown. By convention, a `donothing` project should contain a file called `procedure.md` with the
most recent rendering of this Markdown.

## Default CLI

Since most `donothing` scripts will have the same basic interface, there is a [default
CLI](handle_args.go) that you can use like so:

```
func main() {
    pcd := donothing.NewProcedure()
    // ... set up procedure and steps ...

    // This function parses arguments and performs the appropriate action.
    // Pass --help for a usage message.
    donothing.HandleArgs(os.Args[:], pcd, "root")
}
```
