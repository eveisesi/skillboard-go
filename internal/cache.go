package internal

import (
	"sync"

	"github.com/gofrs/uuid"
)

/**
	The following was created for a very specific purpose.
	I need to be able to propagate a value that is determined in a
	graphql resolver backup up to a middleware that is invoked by the http server.
	I attempted to use context before I realized context is not a pointer
	so I can't add a value to it at the bottom of the chain and
	expected a new/updated value at the top of the chain.
	So I wrote this. It is used in a similar fashion to how I was attempting to use context.
	I'll generate a UUID Key, create an entry in this map with a value of uuid.Nil.
	The UUID Key will be set on the context, passed down and into GraphQL.
	If the intended query is called, then the Resolver will look for this Key in context, do its thing
	and then update this map with the UUID I actually need.
	The middleware is will be sitting there waiting for the response, and when the response is returned
	We'll check this map. If the value of the "Session UUID" is not uuid.Nil, another action will be taken
	at the middleware level and then the entry in the map will be deleted. Simple as that. I almost
	used Redis for this shit since we already use it so heavily in this app. Also, if you find this comment
	and read this far, congrats you get a Gold Star (and maybe ISK if you still play EVE),
	ping me on Tweetfleet Slack (@DoubleD)
**/

var mx = new(sync.Mutex)
var cache = make(map[uuid.UUID]uuid.UUID)

func CacheSet(k uuid.UUID, v uuid.UUID) {
	mx.Lock()
	defer mx.Unlock()
	cache[k] = v
}

func CacheGet(k uuid.UUID) uuid.UUID {
	mx.Lock()
	defer mx.Unlock()
	return cache[k]
}

func CacheDelete(k uuid.UUID) {
	mx.Lock()
	defer mx.Unlock()
	delete(cache, k)
}
