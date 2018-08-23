package junebug

const (
	prompt = `
{
    "text": "JuneBug here to help you setup your meeting!",
    "attachments": [
        {
            "text": "Need video?",
            "fallback": "You will not have any remote attendance",
            "callback_id": "foo",
            "color": "#3AA3E3",
            "attachment_type": "default",
            "actions": [
                {
                    "name": "video",
                    "text": "Jitsi Meet",
                    "type": "button",
                    "value": "jitsi"
                },
                {
                    "name": "video",
                    "text": "Bluejeans",
                    "type": "button",
                    "value": "bluejeans",
                    "confirm" : {
                        "title": "Support for Bluejeans coming :soon:",
                        "ok_text": "So, sorry",
                        "dismiss_text": "So, sorry"
                    }
                }
            ]
        }
    ]
}
`
)
