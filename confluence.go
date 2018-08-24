package junebug

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

// ConfluenceClient generates conference tokens for auth'ed users.
type ConfluenceClient struct {
    URL string
}

// ConfluencePageLinks contains URLs for the created page
type ConfluencePageLinks struct {
    EditUI string `json:"editui"`
    WebUI  string `json:"webui"`
    Self   string `json:"self"`
    Base   string `json:"base"`
}

// ConfluencePage is details about the created page
type ConfluencePage struct {
    ID     string              `json:"id"`
    Type   string              `json:"type"`
    Status string              `json:"status"`
    Title  string              `json:"title"`
    Links  ConfluencePageLinks `json:"_links"`
}

// GetEditURL returns the edit URL for a confluence page
func (cp ConfluencePage) GetEditURL() string {
    return fmt.Sprintf("%s%s", cp.Links.Base, cp.Links.EditUI)
}

// GetWebURL returns the web URL for a confluence page
func (cp ConfluencePage) GetWebURL() string {
    return fmt.Sprintf("%s%s", cp.Links.Base, cp.Links.WebUI)
}

// CreatePage generates a confluence page based on a bunch of stuffs
func (cc ConfluenceClient) CreatePage(token, title, space string) (ConfluencePage, error) {
    createURL := cc.URL + "/rest/api/content"
    jsonStr := []byte(fmt.Sprintf(`{"type":"page","title":"%s","space":{"key":"%s"}}`, title, space))

    req, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonStr))
    if err != nil {
        return ConfluencePage{}, err
    }

    req.Header.Set("Authorization", fmt.Sprintf("JWT %s", token))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return ConfluencePage{}, err
        // handle error
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return ConfluencePage{}, err
        // handle error
    }
    if resp.StatusCode != 200 {
        return ConfluencePage{}, fmt.Errorf("Bad response: %s", body)
        // handle error
    }

    var cp ConfluencePage
    err = json.Unmarshal(body, &cp)
    //    err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&cp)
    if err != nil {
        return ConfluencePage{}, err
    }

    return cp, nil
}
