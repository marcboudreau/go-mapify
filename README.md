go-mapify
=========

This library provides convenience functions to construct Go maps from other
structures.

Usage
-----

Let's say you are using an API that returns a slice of objects but really wished
that you could get them back in a map. You can use mapify to create that map,
with the **FromSlice** function.

```go
package example

import (
    "github.com/marcboudreau/mapify"

    "github.com/third_party/external/api"
)

func disableUserByID(c *api.Client, ids []int) error {
    var users []api.User

    users = c.GetUsers()

    userMap := mapify.FromSlice(users, func(u api.User) int {
        return u.ID()
    })

    for _, id := range ids {
        user, ok := userMap[id]
        if !ok {
            return fmt.Errorf("could not find user with id %d", id)
        }

        user.Disable()
    }

    return nil
}

```

What if you have a use case that prevents using unique ID values for your slice elements? No problem, use the **FromSliceWithDuplicates** function instead.

```go
package example

import (
    "github.com/marcboudreau/go-mapify"

    "github.com/third_party/external/api"
)

func disableUsersInTeam(c *api.Client, teamID int) error {
    var users []api.User

    users = c.GetUsers()

    userMap := mapify.FromSliceWithDuplicates(users, func(u api.User) int {
        return u.TeamID()
    })

    teamUsers, ok := userMap[teamID]
    if !ok {
        return fmt.Errorf("could not find any users for team ID %d", teamID)
    }

    for _, teamUser := range teamUsers {
        teamUser.Disable()
    }

    return nil
}

```