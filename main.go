package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/getlantern/systray"
	"github.com/goodsign/monday"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
	semVer    string // the version of the build
	repoURL   = "https://github.com/MagischerStock/go-calendarweektray" // Repository URL
	apiURL    = "https://api.github.com/repos/MagischerStock/go-calendarweektray/tags" // API URL for tags
)

func main() {
	enableDpiAwareness()

	systray.Run(onReady, func() {})
}

var menus []*systray.MenuItem

func onReady() {
	systray.SetTitle("Kalenderwoche")

	const numberOfEntries = 15
	addMenuItemsForUpcomingCalendarWeekDates(numberOfEntries)

	go keepWeekNumberIconUpToDate()
	go quitOnMenu()
	go checkForUpdates() // Start checking for updates
}

func addMenuItemsForUpcomingCalendarWeekDates(numberOfEntries int) {
	for i := 0; i < numberOfEntries; i++ {
		index := i
		menus = append(menus, systray.AddMenuItem("refresh", ""))
		go func() {
			for {
				<-menus[index].ClickedCh
				_, dateToGo := offsetCalendarWeekToDate(index)
				goToDate(dateToGo)
			}
		}()
	}

	systray.AddSeparator()
}

func keepWeekNumberIconUpToDate() {
	calendarWeek := currentCalendarWeekIterator()
	for {
		updateIconAndTooltip(<-calendarWeek.ChangedCh)
	}
}

func quitOnMenu() {
	quitMenuItem := systray.AddMenuItem(fmt.Sprintf("Beenden (%s - %s)", semVer, sha1ver), "Beendet die Applikation")
	<-quitMenuItem.ClickedCh
	systray.Quit()
}

func updateIconAndTooltip(weekNumber int) {
	systray.SetIcon(generateImage(weekNumber))
	systray.SetTooltip(fmt.Sprintf("Aktuelle Kalenderwoche: %d", weekNumber))

	refreshUpcomingCalendarWeekItems()
}

func refreshUpcomingCalendarWeekItems() {
	for index, _ := range menus {
		week, startDate := offsetCalendarWeekToDate(index)

		text := fmt.Sprintf("KW %d: %s", week, monday.Format(startDate, "02. January 2006", monday.LocaleDeDE))
		menus[index].SetTitle(text)
	}
}

func checkForUpdates() {
	// Add menu item to show version
	versionMenuItem := systray.AddMenuItem(fmt.Sprintf("Version: %s", semVer), "Aktuelle Version")

	// Fetch latest release tag from GitHub API
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error fetching tags:", err)
		return
	}
	defer resp.Body.Close()

	var tags []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		fmt.Println("Error decoding tags:", err)
		return
	}

	if len(tags) > 0 && tags[0].Name != semVer {
		// Update available
		versionMenuItem.SetTitle(fmt.Sprintf("Version: %s (Update verfügbar!)", semVer))
		go func() {
			for {
				<-versionMenuItem.ClickedCh
				openBrowser(repoURL)
			}
		}()
	} else {
		// No update
		go func() {
			for {
				<-versionMenuItem.ClickedCh
				openBrowser(repoURL)
			}
		}()
	}
}

func openBrowser(url string) {
	// Function to open the URL in the default browser
	fmt.Println("Opening browser with URL:", url)
}
