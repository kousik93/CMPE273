try:
    import requests
    import datetime
    import json
except Exception as e:
    print "Requests library not found. Please install it. \nHint: pip install requests"

person = {

    "email": "foo@gmail.com",
    "zip": "94112",
    "country": "U.S.A",
    "profession": "student",
    "favorite_color": "blue",
    "is_smoking": "yes|no",
    "favorite_sport": "hiking",
    "food": {
        "type": "vegetrian|meat_eater|eat_everything",
        "drink_alcohol": "yes|no"
    },
    "music": {
        "spotify_user_id": "wizzler"
    },
    "movie": {
        "tv_shows": ["x", "y", "z"],
        "movies": ["x", "y", "z"]
    },
    "travel": {
        "flight": {
            "seat": "aisle|window"
        }
    }
}

change_person = {
    "travel": {
        "flight": {
            "seat": "middle"
        }
    },
    "movie": {
        "tv_shows": ["a", "a", "a"],
        "movies": ["a", "a", "a"]
    },
    "food": {
        "type": "vegetrian",
        "drink_alcohol": "yes"
    },
    "favorite_sport": "football"
}


def test_get(url):
    get_url = "%s/%s" % (url, person['email'])
    print get_url
    r = requests.get(get_url)
    print r.status_code
    print r.json()
    if r.status_code == 200:
        d = r.json()
        try:
            d = d[0]
        except KeyError:
            pass
        if d["zip"] == person['zip']:
            if type(d['food']) is dict:
                if d['movie']['movies'][0] == person['movie']['movies'][0]:
                    t = datetime.datetime.now()
                    print "GET check successful. Time: %s" % t
    else:
        print "GET check failed"



def test_post(url):
    post_url = url
    try:
        r = requests.post(post_url, data=json.dumps(person))
        print r.status_code
        if r.status_code == 201:
            d = datetime.datetime.now()
            print "POST Check successful. Time: %s" % d
        else:
            print "POST incorrect. Not working as expected."
    except requests.exceptions.ConnectionError:
        print "Server not running on %s" % url
        exit()



def test_delete(url):
    delete_url = "%s/%s" % (url, person['email'])
    r = requests.delete(delete_url)
    print r.status_code
    if r.status_code == 204:
        print "DELETE status code check complete"
    else:
        print "DELETE failed."

    get_url = "%s/profile/%s" % (url, person['email'])
    r = requests.get(url)
    if r.status_code == 200:
        print "DELETE has not deleted the item."
    else:
        print "You have deleted the item"



def test_put(url):
    put_url = "%s/%s" % (url, person['email'])
    r = requests.put(put_url, data=json.dumps(change_person))
    if r.status_code == 204:
        t = datetime.datetime.now()
        print "PUT request sent successfully. Time: %s" % t
    else:
        print "PUT failed."



test_post("http://localhost:3000/profile")
test_put("http://localhost:3000/profile")
