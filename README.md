# art-detective

An Artsy API inspector made for Eafit's Component-Based Software Development course.

## What?

Hypermedia APIs represent relationships between objects in their responses with a link to the other
resources. **art-detective** is just a CLI client for Artsy's API that allows you to follow those
hyperlinks.

## Example

To follow this example, clone the repository and `cd` into its folder.
```
git clone https://github.com/castillobg/art-detective.git
cd art-detective
```
Then, export your client ID and client secret to environment variables.
```
export ARTSY_CLIENT_ID=<your client id>
export ARTSY_CLIENT_SECRET=<your client secret>
```

By default, art-detective will direct requests to the `artwork` endpoint. So, with no arguments,
it lists many artworks sent back from the api for the request `GET https://api.artsy.net/api/artworks`.

You can see how that looks by running
```
./art-detective
```

You can specify another endpoint, or _subject_ for art-detective to investigate (**artworks** and **artists** are currently the only supported subjects):
```
./art-detective -subject artists
```

If you'd like art-detective to investigate a specific subject, you can pass an id as an argument. For example, let's investigate Andy Warhol, whose ID is 4d8b92b34eb68a1b2c0003f4:
```
./art-detective -subject artists -id 4d8b92b34eb68a1b2c0003f4
```

Because Artsy's is a Hypermedia API, you'll see some fields are URLs pointing to a related resource. Let's take a look at the `_links.similar_artists.href` field:

```
./art-detective -subject artists -id 4d8b92b34eb68a1b2c0003f4 -field _links.similar_artists.href
```

art-detective first gets the object for Andy Warhol, and then sends a request to the URL present in the specified field and prints the response.
