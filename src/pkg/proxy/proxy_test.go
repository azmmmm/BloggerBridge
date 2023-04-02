package proxy

import (
	"net/http"
	"testing"
)

func TestFetchByProxy(t *testing.T) {
	// Replace this with the URL you want to test
	//testURL := "https://blogger.googleusercontent.com/img/a/AVvXsEiddYtOmrnxvAXCNI65yN8ArL9Xv_VgN6JS2WOoI9Zr4092izKbWXuUwTMJafMe0y_6VqRy-Zr9ZYVa70CsQp_P5KLPaRsBp7mJH-Rb7AOuFH2VUPBW1IG2hn2t4Mjvh-aYwDUtwhpuHM-rjD6nFBQLlmR3pnveTAhZj331lbzV0UOG8Xn3kEldiH8klQ=w480"
	//testURL := "https://baidu.com"
	testURLs := []string{
		"https://www.youtube.com",
		"https://baidu.com",
		"https://blogger.googleusercontent.com/img/a/AVvXsEiddYtOmrnxvAXCNI65yN8ArL9Xv_VgN6JS2WOoI9Zr4092izKbWXuUwTMJafMe0y_6VqRy-Zr9ZYVa70CsQp_P5KLPaRsBp7mJH-Rb7AOuFH2VUPBW1IG2hn2t4Mjvh-aYwDUtwhpuHM-rjD6nFBQLlmR3pnveTAhZj331lbzV0UOG8Xn3kEldiH8klQ=w480",
	}
	for _, url := range testURLs {
		testURL(t, url)
	}
}

func testURL(t *testing.T, url string) {
	resp, err := FetchByProxy(url)
	if err != nil {
		t.Fatalf("FetchByProxy failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
