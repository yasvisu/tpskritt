//TPSkritt is an application of the Go GW2API: https://github.com/yasvisu/gw2api .
//
//To use this package, run it with go run.
//
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/yasvisu/gw2api"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

//CONSTANTS

const help = "/q to quit\n" +
	"help for help\n" +
	"/save to save your settings and subscriptions\n" +
	"/sub {item} to subscribe for the item\n" +
	"/subid to subscribe ids in bulk\n" +
	"/subs for a list of subscriptions\n" +
	"/d /del /delete to delete subscriptions\n" +
	"/flush to flush all subscriptions\n" +
	"/nosave to quit without saving\n" +
	"\n"+
	"/user to set your username\n" +
	"enter to see all subscription prices\n"
	
var user string
const sys = "tpskritt"
const ver = "v0.1"
const MAXSUBS = 20

//INIT

func init() {
	blob, err := ioutil.ReadFile("settings.ini")
	if err == nil {
		user = string(blob)
	} else {
		user = "user"
	}
	chatter("initializing!")
	subsNames = make(map[int]string)
	
	//GET SUBSCRIPTIONS AND FLAVOR
	blob, err = ioutil.ReadFile("subscriptions.ini")
	if err == nil {
		res := strings.Split(string(blob), "\n")
		for _, val := range res[:len(res)-1] {
			tmp, err := strconv.Atoi(val)
			if err != nil {
				chatter("subscription id entry cannot be parsed!\n" + val)
				continue
			}
			subs = append(subs, tmp)
		}
	}
	
	items, _ := gw2api.ItemsIds("", subs...)
	for _, val := range items {
		subsNames[val.ID] = val.Name
	}
}

//MAIN

func main() {
	fmt.Println("Welcome to tpskritt " + ver + "!")
	var wg sync.WaitGroup
	wg.Add(1)
	go inputDaemon(&wg)
	wg.Wait()
	os.Exit(0)
}

//INPUT

func inputDaemon(wg *sync.WaitGroup) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", user)
		text, err := reader.ReadString('\n')
		text = text[:len(text)-2]

		if err != nil {
			fmt.Println(err)
		}

		res := strings.Fields(text)

		if len(res) > 0 {
			switch res[0] {
			case "/q":
				saveSubscriptions()
				saveSettings()
				chatter("quitting")
				wg.Done()
				return
			case "/nosave":
				chatter("quitting without saving")
				wg.Done()
				return
			case "/sub":
				//subscribe item
				//empty sub with no previous
				if len(text) > 5 {
					subscribe(text[5:])
				} else {
					listSubscriptions()
				}
			case "/subid":
				var tmp []int
				for _, val := range strings.Split(text, " ")[1:] {
					t, err := strconv.Atoi(val)
					if err != nil {
						continue
					}
					tmp = append(tmp, t)
				}
				subscribeAll(tmp)
			case "/subs":
				//show subscriptions
				listSubscriptions()
			case "/d":
				fallthrough
			case "/del":
				fallthrough
			case "/delete":
				var tmp []int
				for _, val := range strings.Split(text, " ")[1:] {
					t, err := strconv.Atoi(val)
					if err != nil {
						continue
					}
					tmp = append(tmp, t)
				}
				deleteSubs(tmp)
			case "/s":
				fallthrough
			case "/save":
				saveSubscriptions()
				saveSettings()
			case "/user":
				user = text[6:]
			case "/flush":
				subs = nil
				subsNames = make(map[int]string)
				chatter("flushed all subscriptions")
			case "help":
				fallthrough
			case "-h":
				fallthrough
			case "/h":
				fmt.Println(help)
			default:
				query(text)
			}
		} else {
			analyzeSubscriptions()
		}
	}
}

//Look up query string and return a slice of GW2API article prices.
func lookup(q string) (p []gw2api.ArticlePrices) {
	num, err := strconv.Atoi(q)
	if err == nil {
		res, err := gw2api.CommercePricesIds(num)
		if err != nil {
			chatter(err.Error())
			return nil
		}
		return res
	} else {
		spidyItemIDs, err := spidySearchID(q)
		if err != nil {
			chatter(err.Error())
			return nil
		}
		if len(spidyItemIDs) < 1 {
			chatter("No results")
			return nil
		} else if len(spidyItemIDs) > 1 {
			for _, val := range spidyItemIDs {
				fmt.Printf("\t⇄%d\t%s\n", val.ID, val.Name)
			}
			return nil
		} else {
			res, err := gw2api.CommercePricesIds(spidyItemIDs[0].ID)
			if err != nil {
				chatter(err.Error())
				return
			}
			return res
		}
	}
}

//Look up and analyze a query.
func query(q string) {
	prices := lookup(q)
	if prices == nil {
		return
	}
	analyze(prices[0])
}

//ANALYSIS
//INCOMPLETE
/*
type Analysis struct {
}
*/

//LISTING FEE: 5% on
//TRANSACTION FEE: 10%
func analyze(p gw2api.ArticlePrices) {
	i, err := gw2api.ItemsIds("", p.ID)
	if err != nil {
		chatter(err.Error())
	}
	fmt.Printf("\t%s\t%s(%d)\t%s(%d)\n",
		i[0].Name,
		toCurrency(p.Buys.UnitPrice), p.Buys.Quantity,
		toCurrency(p.Sells.UnitPrice), p.Sells.Quantity)

	/*analyzeOutput(p)
	analyzeInput(p)*/
	/*recipeO, err := gw2api.RecipesSearchOutput("", p.ID)
	recipeI, err := gw2api.RecipesSearchInput("", p.ID)
	recipes, err := gw2api.RecipesIds(recipeO..., recipeI...)*/
}

//Analyze all subscriptions.
func analyzeSubscriptions() {
	items, err := gw2api.ItemsIds("", subs...)
	if err != nil {
		chatter(err.Error())
		return
	}
	prices, err := gw2api.CommercePricesIds(subs...)
	if err != nil {
		chatter(err.Error())
		return
	}
	if len(items) != len(items) {
		chatter("Mismatching items list size and prices list size")
		return
	}
	for key, _ := range items {
		fmt.Printf("\t%d. %s\t%s(%d)\t%s(%d)\n",
			key,
			items[key].Name,
			toCurrency(prices[key].Buys.UnitPrice), prices[key].Buys.Quantity,
			toCurrency(prices[key].Sells.UnitPrice), prices[key].Sells.Quantity)
	}
}

//SUBSCRIPTIONS

//sorted list of subscription ids.
var subs []int

//map of id - name for each subscription.
var subsNames map[int]string

//Subscribe item (if it is a valid one) and add it to the map.
func subscribe(q string) {
	items := lookup(q)
	if len(items) != 1 {
		return
	} else if contains(subs, items[0].ID) {
		chatter("already subscribed")
		return
	}

	subs = append(subs, items[0].ID)
	
	tmp, _ := gw2api.ItemsIds("",items[0].ID)
	subsNames[items[0].ID] = tmp[0].Name
	sort.Ints(subs)
	chatter("subscribed")
}

//Subscribe all ids.
func subscribeAll(is []int) {
	items, err := gw2api.ItemsIds("", is...)
	if err != nil {
		chatter(err.Error())
		return
	}
	for _, val := range items {
		if contains(subs, val.ID) {
			chatter("already subscribed: " + strconv.Itoa(val.ID))
		} else {
			subs = append(subs, val.ID)
			subsNames[val.ID] = val.Name
			chatter(strconv.Itoa(val.ID) + " subscribed OK")
		}
	}
	sort.Ints(subs)
	chatter("done subscribing " + strconv.Itoa(len(items)) + "/" + strconv.Itoa(len(is)) + " ids")
}

//Checks whether the slice A contains the element B.
func contains(a []int, b int) bool {
	for _, val := range a {
		if val == b {
			return true
		}
	}
	return false
}

//List all subscription ids with their corresponding names.
func listSubscriptions() {
	if len(subs) == 0 {
		return
	}
	for key, val := range subs {
		fmt.Printf("\t%d⇄%d\t%s\n", key, val, subsNames[val])
	}
}

//Save all subscriptions to subscription.ini.
func saveSubscriptions() {
	file, err := os.Create("subscriptions.ini")
	defer file.Close()
	if err != nil {
		chatter("Could not save: creating subscriptions.ini failed")
		return
	}
	if len(subs) > MAXSUBS {
		subs = subs[len(subs)-(MAXSUBS-1):]
	}
	for _, val := range subs {
		file.WriteString(strconv.Itoa(val) + "\n")
	}
	chatter("saved subscriptions!")
}

//Delete subscriptions by their index.
func deleteSubs(is []int) {
	sort.Ints(is)
	str := ""
	for key, val := range is {
		if val - key >= 0 && val - key < len(subs) {
			copy(subs[val-key:], subs[val-key+1:])
			subs[len(subs)-1] = 0
			subs = subs[:len(subs)-1]
			str += strconv.Itoa(val) + " "
		}
	}
	chatter("deleted " + str[:len(str)-1] + "!")
	listSubscriptions()
}

//Save all settings to settings.ini.
func saveSettings() {
	file, err := os.Create("settings.ini")
	defer file.Close()
	if err != nil {
		chatter("Could not save: creating settings.ini failed")
		return
	}
	file.WriteString(user)
	chatter("saved settings!")
}

//SPIDY

type idName struct {
	ID   int
	Name string
}

func spidySearchID(q string) (res []idName, err error) {
	type spidyItem struct {
		DataID int    `json:"data_id"`
		Name   string `json:"name"`
	}
	type spidyData struct {
		Results []spidyItem `json:"results"`
	}

	resp, err := http.Get("http://www.gw2spidy.com/api/v0.9/json/item-search/" + q)
	if err != nil {
		return nil, err
	}

	blob, err := ioutil.ReadAll(resp.Body)

	var spd spidyData
	err = json.Unmarshal(blob, &spd)

	for _, val := range spd.Results {
		var t idName
		t.ID = val.DataID
		t.Name = val.Name
		res = append(res, t)
	}

	return
}

//MISC

func chatter(t string) {
	fmt.Printf("%s: %s\n", sys, t)
}

func toCurrency(i int) (res string) {
	res += strconv.Itoa(i/10000) + "g "
	i = i % 10000
	res += strconv.Itoa(i/100) + "s "
	i = i % 100
	res += strconv.Itoa(i) + "c"
	return
}
