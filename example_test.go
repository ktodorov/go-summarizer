package goSummarizer

import (
	"fmt"
)

func ExampleCreateFromText() {
	var unsummarizedText = "unsummarized text"
	var s = CreateFromText(unsummarizedText)
	// Do something with s
	s.Summarize()
}

func ExampleCreateFromURL() {
	var urlToSummarize = "http://testurl.test/"
	var s = CreateFromURL(urlToSummarize)
	// Do something with s
	s.Summarize()
}

func ExampleSummarizer_Summarize() {
	var customNewsStory = `SpaceX has succeeded in launch a Falcon 9 rocket from Vandenberg Air Force Base in California, its first launch since a Falcon 9 rocket exploded on a launch pad in pre-flight procedures in September 2016. The launch took place at 9:54 AM PT Saturday, during an instant launch window. 

This mission is the first in a series for client Iridium, that will see it deploy 70 satellites in a network for voice and data communication. It’s also a green light for SpaceX in terms of the company pursuing its aggressive launch schedule, which is something the private launch provider needs to do in order to continue locking in new contracts and working towards its goal of decreasing the cost of launches even further still.

In 2016, SpaceX completed only 8 of a planned 20 launches, due to the September 1 explosion that halted all new launches for four months. That has not been good for the company’s bottom line, resulting in a year that likely saw it exacerbate a reported $250 million loss in 2015.
SpaceX also had to push back its timelines for test launches of its Dragon crew capsule as a result of the September incident. The original target date for a Dragon test launch with people on board was 2017, but it’s now been pushed back to 2018. The company still hopes to fly a mission without crew on board by the last quarter of this year, however.
Crewed mission capabilities will help SpaceX expand its ability to serve contracts, since it can then serve the ISS for more than just supply runs. It also sets the stage for SpaceX’s future goals of providing missions to Mars, with a target initial date for those aspirations still set for 2024.`

	var s = CreateFromText(customNewsStory)
	summary, err := s.Summarize()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(summary)
	// Output: SpaceX has succeeded in launch a Falcon 9 rocket from Vandenberg Air Force Base in California, its first launch since a Falcon 9 rocket exploded on a launch pad in pre-flight procedures in September 2016
	// It’s also a green light for SpaceX in terms of the company pursuing its aggressive launch schedule, which is something the private launch provider needs to do in order to continue locking in new contracts and working towards its goal of decreasing the cost of launches even further still.
	// SpaceX also had to push back its timelines for test launches of its Dragon crew capsule as a result of the September incident
}

func ExampleSummarizer_Summarize_second() {
	var customNewsStoryURL = `https://techcrunch.com/2017/01/14/spacex-successfully-returns-to-launch-with-iridium-1-next-falcon-9-mission/`

	var s = CreateFromURL(customNewsStoryURL)
	summary, err := s.Summarize()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(summary)
	// Output: The launch also resulted in a successful recovery of the Falcon 9 rocket’s first stage, which marks the seventh time SpaceX has succeed in landing this stage back for potential later re-use
	// SpaceX also had to push back its timelines for test launches of its Dragon crew capsule as a result of the September incident
	// All satellites were successfully deployed as of 11:13 AM PT / 2:12 PM PT, signalling a successful mission for the space company’s first flight back.
}

func ExampleSummarizer_GetSummaryInfo() {
	var s = CreateFromText("test")
	s.Summarize()
	summaryInfo, err := s.GetSummaryInfo()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(summaryInfo)
	// Output: Summary info:
	//  - Original length: 4
	//  - Summary length:  0
	//  - Summary ratio:   100
}

func ExampleSummarizer_IsSummarized() {
	var s = CreateFromText("test")
	fmt.Println(s.IsSummarized())
	// Output: false
}

func ExampleSummarizer_IsSummarized_second() {
	var s = CreateFromText("test")
	s.Summarize()
	fmt.Println(s.IsSummarized())
	// Output: true
}

func ExampleSummarizer_StoreToFile() {
	var s = CreateFromText("test")
	s.Summarize()
	stored, err := s.StoreToFile("some/path/to/file.txt")
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(stored)
	// Output: true
}
