# go-summarizer
This is a Go library for summarizing text and websites and optionally saving the data to a local file

[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)

## Installing
    go get github.com/ktodorov/go-summarizer

## Creating Summarizer instance

### From text
    var unsummarizedText = "unsummarized text"
	var s = CreateFromText(unsummarizedText)

### From website url
    var urlToSummarize = "http://testurl.test/"
	var s = CreateFromURL(urlToSummarize)

## Supported methods
### Summarize
    var customNewsStoryURL = `https://techcrunch.com/2017/01/14/spacex-successfully-returns-to-launch-with-iridium-1-next-falcon-9-mission/`
	
    var s = CreateFromURL(customNewsStoryURL)
	summary, err := s.Summarize()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
        return
	}

	fmt.Println(summary)

Output*:
>SpaceX successfully returns to launch with Iridium-1 NEXT Falcon 9 mission
>
>It’s a huge victory for SpaceX, which has had to delay its launch schedule since the explosion.
>The launch also resulted in a successful recovery of the Falcon 9 rocket’s first stage, which marks the seventh time SpaceX has succeed in landing this stage back for potential later re-use
>It’s also a green light for SpaceX in terms of the company pursuing its aggressive launch schedule, which is something the private launch provider needs to do in order to continue locking in new contracts and working towards its goal of decreasing the cost of launches even further still.
>In 2016, SpaceX completed only 8 of a planned 20 launches, due to the September 1 explosion that halted all new launches for four months
>SpaceX also had to push back its timelines for test launches of its Dragon crew capsule as a result of the September incident
>It also sets the stage for SpaceX’s future goals of providing missions to Mars, with a target initial date for those aspirations still set for 2024.
>All satellites were successfully deployed as of 11:13 AM PT / 2:12 PM PT, signalling a successful mission for the space company’s first flight back.

_*Note that it first prints the title of the web page if there is such_

### GetSummaryInfo
    var s = CreateFromText("first sentence. second sentence")
	s.Summarize()
	summaryInfo, err := s.GetSummaryInfo()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(summaryInfo)

Output:
>Summary info: <br/>
> \- Original length: 31 symbols <br/>
> \- Summary length:  14 symbols <br/>
> \- Summary ratio:   54.84% <br/>

### IsSummarized
    var s = CreateFromText("first sentence. second sentence")
	fmt.Println("Before summarizing: ", s.IsSummarized())
	s.Summarize()
	fmt.Println("After summarizing: ", s.IsSummarized())

Output:
> Before summarizing: false<br/>
> After summarizing: true<br/>

### StoreToFile
    var s = CreateFromText("first sentence. second sentence")
	s.Summarize()
	stored, err := s.StoreToFile("some/path/to/file.txt")
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(stored)

Output:
> true

_*Currently supported file types: txt and pdf_