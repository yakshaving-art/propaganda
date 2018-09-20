package webhook

// Headers
// "X-Hub-Signature":[]string{"sha1=be490c94029284a1074f6ed7d6f551affcfa6e8b"},
// "User-Agent":[]string{"GitHub-Hookshot/32d792e"},
// "Content-Length":[]string{"20774"},
// "X-Github-Delivery":[]string{"97757e90-bb8e-11e8-8017-464837e8ed07"},
// "X-Github-Event":[]string{"pull_request"},
// "Accept-Encoding":[]string{"gzip"},
// "Accept":[]string{"*/*"},
// "Content-Type":[]string{"application/json"}

// // Repository holds the repository information
// type Repository struct {
// 	URL      string `json:"url"`
// 	FullName string `json:"full_name"`
// }

// // Hook holds the hook events
// type Hook struct {
// 	Events []string `json:"events"`
// }

// // HookPayload hold the GitHub
// type HookPayload struct {
// 	Repository Repository `json:"repository"`
// 	Hook       Hook       `json:"hook"`
// }

// // GetRepository implements webhook.HookPayload interface
// func (h HookPayload) GetRepository() string {
// 	return h.Repository.FullName
// }

// // ParseHookPayload parses a payload string and returns the payload as a struct
// func (c Client) ParseHookPayload(payload string) (webhooks.HookPayload, error) {
// 	var hookPayload HookPayload
// 	if err := json.Unmarshal([]byte(payload), &hookPayload); err != nil {
// 		return hookPayload, fmt.Errorf("could not parse hook payload: %s", err)
// 	}
// 	return hookPayload, nil
// }
