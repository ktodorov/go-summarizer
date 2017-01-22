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

Output:
>The launch also resulted in a successful recovery of the Falcon 9 rocket’s first stage, which marks the seventh time SpaceX has succeed in landing this stage back for potential later re-use<br/>
>SpaceX also had to push back its timelines for test launches of its Dragon crew capsule as a result of the September incident<br/>
>All satellites were successfully deployed as of 11:13 AM PT / 2:12 PM PT, signalling a successful mission for the space company’s first flight back.<br/>

### GetSummaryInfo
    var s = CreateFromText("test")
	s.Summarize()
	summaryInfo, err := s.GetSummaryInfo()
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(summaryInfo)

Output:
>Summary info:<br/>
> \- Original length: 4<br/>
> \- Summary length:  0<br/>
> \- Summary ratio:   100<br/>

### IsSummarized
    var s = CreateFromText("test")
	fmt.Println("Before summarizing: ", s.IsSummarized())
	s.Summarize()
	fmt.Println("After summarizing: ", s.IsSummarized())

Output:
> Before summarizing: false<br/>
> After summarizing: true<br/>

### StoreToFile
    var s = CreateFromText("test")
	s.Summarize()
	stored, err := s.StoreToFile("some/path/to/file.txt")
	if err != nil {
		fmt.Println("Error occurred: ", err.Error())
	}

	fmt.Println(stored)

Output:
> true