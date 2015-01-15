# tpskritt

##[Guild Wars 2](https://www.guildwars2.com/en-gb/) - Trading Post tool written on [Go](http://golang.org/).
This tool is also meant to also be an application of the [Go GW2API wrapper](https://github.com/yasvisu/gw2api).

Note:
This is an experimental command-line tool and may change in functionality at any moment.

## Requirements

* Developed with [Go](http://golang.org/) 1.3.3. 

## Installation
Make sure you have [Go installed and set up](http://golang.org/doc/install)!
To get the tpskritt package, simply hit in your console of choice:

    go get github.com/yasvisu/tpskritt

In your favorite console, navigate to your github package directory and have a look at the tpskritt.go code.
To open a console on Windows, hit:

    ctrl+R

Type in cmd, hit enter, and then navigate with `cd ...` to the directory of tpskritt.go.
If you like everything you see in the code, you can run tpskritt in two ways:

Run the code directly:

    go run tpskritt.go

Install the code and run it as a binary:

    go install tpskritt.go

After installing, find it in your bin (GOBIN) folder and run it. You may want to move it to a separate directory, as it creates .ini files to store its settings and subscriptions.

## Documentation
What you can do with tpskritt:
- look up items from the GW2Spidy database (thanks to GW2Spidy's API!), based on ID or name.
- look up item prices from the GW2 commerce database.
- subscribe to multiple items, so you can track their trading post prices with the stroke of a button

### Commands:
  `/q` to quit
  
  `help` for help
  
  `/save` `/s` to save your settings and subscriptions
  
  `/sub` {item} to subscribe for the item
  
  `/subid` to subscribe ids in bulk
  
  `/subs` for a list of subscriptions
  
  `/d` `/del` `/delete` to delete subscriptions
  
  `/flush` to flush all subscriptions
  
  `/nosave` to quit without saving
  
  `/user` to set your username

## Development state
###V0.1
This is an experimental release which I have made for fun. There's a huge to-do list here, so I'm stopping work on this, unless there's requests for me to continue.
* Experimental release
* Polishing and debugging ongoing by default
* To do:
  * Get feedback
  * Add flavor
      * sys names
      * skritt dialogue
  * Add features
      * finish experimental analysis
