package models

import (
	"errors"
	"log"
	"strconv"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

// Declare a global pl variable to store the Redis connection pool.
var pl *pool.Pool

func init() {
	var err error
	// Establish a pool of 10 connections to the Redis server listening on
	// port 6379 of the local machine.
	pl, err = pool.New("tcp", "localhost:6379", 10)
	if err != nil {
		log.Println(err)
	}
}

// Create a new error message and store it as a constant. We'll use this
// error later if the FindAlbum() function fails to find an album with a
// specific id.

var ErrNoAlbum = errors.New("models: no album found")

type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func FindTopThree() ([]*Album, error) {
	conn, err := pl.Get()
	if err != nil {
		return nil, err
	}
	defer pl.Put(conn)

	// Begin an infinite loop.

	for {
		// Instruct Redis to watch the likes sorted set for any changes.

		err = conn.Cmd("WATCH", "likes").Err
		if err != nil {
			return nil, err
		}
		// Use the ZREVRANGE command to fetch the album ids with the highest
		// score (i.e. most likes) from our 'likes' sorted set. The ZREVRANGE
		// start and stop values are zero-based indexes, so we use 0 and 2
		// respectively to limit the reply to the top three. Because ZREVRANGE
		// returns an array response, we use the List() helper function to
		// convert the reply into a []string.
		reply, err := conn.Cmd("ZREVRANGE", "likes", 0, 2).List()
		if err != nil {
			return nil, err
		}

		// Use the MULTI command to inform Redis that we are starting a new
		// transaction.
		for _, id := range reply {
			err := conn.Cmd("HGETALL", "album: "+id).Err
			if err != nil {
				return nil, err
			}
		}

		// Execute the transaction. Importantly, use the Resp.IsType() method
		// to check whether the reply from EXEC was nil or not. If it is nil
		// it means that another client changed the WATCHed likes sorted set,
		// so we use the continue command to re-run the loop.
		ereply := conn.Cmd("EXEC")
		if ereply.Err != nil {
			return nil, err
		} else if ereply.IsType(redis.Nil) {
			continue
		}
		// Otherwise, use the Array() helper function to convert the
		// transaction reply to an array of Resp objects ([]*Resp).
		areply, err := ereply.Array()
		if err != nil {
			return nil, err
		}

		// Create a new slice to store the album details.

		abs := make([]*Album, 3)
		// Iterate through the array of Resp objects, using the Map() helper
		// to convert the individual reply into a map[string]string, and then
		// the populateAlbum function to create a new Album object
		// from the map. Finally store them in order in the abs slice.
		for i, reply := range areply {
			mreply, err := reply.Map()
			if err != nil {
				return nil, err
			}
			ab, err := populateAlbum(mreply)
			if err != nil {
				return nil, err
			}
			abs[i] = ab
		}

		return abs, nil
	}
}

func IncrementLikes(id string) error {
	conn, err := pl.Get()
	if err != nil {
		return err
	}
	defer pl.Put(conn)
	// Before we do anything else, check that an album with the given id
	// exists. The EXISTS command returns 1 if a specific key exists
	// in the database, and 0 if it doesn't.
	exists, err := conn.Cmd("EXISTS", "album: "+id).Int()
	if err != nil {
		return err
	} else if exists == 0 {
		return ErrNoAlbum
	}

	// Use the MULTI command to inform Redis that we are starting a new
	// transaction.
	err = conn.Cmd("MULTI").Err
	if err != nil {
		return err
	}

	// Increment the number of likes in the album hash by 1. Because it
	// follows a MULTI command, this HINCRBY command is NOT executed but
	// it is QUEUED as part of the transaction. We still need to check
	// the reply's Err field at this point in case there was a problem
	// queueing the command.

	err = conn.Cmd("HINCRBY", "album:"+id, "likes", 1).Err
	if err != nil {
		return err
	}
	// And we do the same with the increment on our sorted set.
	err = conn.Cmd("ZINCRBY", "likes", 1, id).Err
	if err != nil {
		return err
	}
	// Execute both commands in our transaction together as an atomic group.
	// EXEC returns the replies from both commands as an array reply but,
	// because we're not interested in either reply in this example, it
	// suffices to simply check the reply's Err field for any errors.
	err = conn.Cmd("EXEC").Err
	if err != nil {
		return err
	}
	return nil
}

func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	ab := new(Album)
	ab.Title = reply["title"]
	ab.Artist = reply["artist"]
	ab.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}
	ab.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		return nil, err
	}
	return ab, nil
}

func FindAlbum(id string) (*Album, error) {
	// Use the connection pool's Get() method to fetch a single Redis
	// connection from the pool.
	// conn, err := pl.Get()
	// if err != nil {
	//  return nil, err
	// }
	// Importantly, use defer and the connection pool's Put() method to ensure
	// that the connection is always put back in the pool before FindAlbum()
	// exits.
	// defer pl.Put(conn)

	// Fetch the details of a specific album. If no album is found with the
	// given id, the map[string]string returned by the Map() helper method
	// will be empty. So we can simply check whether it's length is zero and
	// return an ErrNoAlbum message if necessary.
	reply, err := pl.Cmd("HGETALL", "album:"+id).Map()
	if err != nil {
		return nil, err
	} else if len(reply) == 0 {
		return nil, ErrNoAlbum
	}

	return populateAlbum(reply)
}
