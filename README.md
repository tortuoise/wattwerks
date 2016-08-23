<A name="toc0_1" title="Wattwerks"/>
# Example crud catalog in Go on Appengine

##Contents     
**<a href="toc1_1">Handlers</a>**  
**<a href="toc1_2">Datastore</a>**  
**<a href="toc1_3">Templates</a>**  
**<a href="toc1_4">Testing</a>**   


<A name="toc1_1" title="Handlers" />
## Handlers ##
All the handlers are currently in wattwerks.go. TODO: needs some tidying
<A name="toc1_2" title="Datastore" />
## Datastore ##
The datastore provides read and write access to the appengine datastore. Main guideline to follow is that no datastore data structure should ever be used outside this package. Another guideline is to have all functions from this package throw DSErr which implements error interface.
<A name="toc1_3" title="Templates" />
## Templates ##
All the templates are in wattwerks.go. TODO: needs some tidying
<A name="toc1_4" title="Tests" />
## Testing ##
Testing with the dev_appserver needs environment variable $APPENGINE_DEV_APPSERVER set. Run tests as usual with go test ./wattwerks
