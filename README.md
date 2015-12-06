GitHub webhook 
===============

A simple webhook in [Go](https://golang.org/) for GitHub to check that any poll requst contains "Signed-off-by:" line.


HowTo use deployed instance
============================
   - go to repositories Settings->Webhooks & services
   - add new web hook
   - put https://signed-off-by-verifier.herokuapp.com as Payload URL
   - required permissins 'pull_requests" only
   - done! 
   


