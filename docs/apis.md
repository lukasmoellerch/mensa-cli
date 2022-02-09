# Mensa APIs

There are several APIs which can be used to fetch data about currently available meals at ETH and UZH mensas. I didn't find any publicly available documentation about them, but I'll do my best to collect my findings here.

## ETHZ Webservices

The gastronomy API is available at `www.webservices.ethz.ch/gastro/v1`. I found the following endpoints:

### `www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/facilities/`

Returns a list of all gastronomy facilities.

```
www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/facilities/
```
```json
{
  "locations": [
    {
      "id": 1,
      "label": "Zentrum",
      "label_en": "Zentrum"
    }
  ],
  "facilites": [
    {
      "id": 2,
      "label": "bQm",
      "label_en": "bQm",
      "type": "bar",
      "location_id": 1
    },
    {
      "id": 4,
      "label": "Clausiusbar",
      "label_en": "Clausiusbar",
      "type": "mensa",
      "location_id": 1
    }
  ]
}
```

### `www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/facilities/{id}/{language}`

Returns details about the specified facility.

```
www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/facilities/12/de
```
```json
{
  "id": 12,
  "caterer": null,
  "type": "mensa",
  "label": "Mensa Polyterrasse",
  "internal_link": "https://www.ethz.ch/{lang}/campus/erleben/gastronomie-und-einkaufen/gastronomie/restaurants-und-cafeterias/zentrum/mensa-polyterrasse.html",
  "features": [
    {
      "id": 4,
      "label": "warme Hauptmahlzeiten"
    }
  ],
  "hours": [
    {
      "date": "03.01.2022",
      "day": "Montag",
      "hours": {
        "opening": [],
        "mealtime": []
      }
    }
  ],
  "address": {
    "building": "MM",
    "room": "B 79",
    "street": "Leonhardstrasse 34",
    "city": "8092 Zürich"
  },
  "managers": [
    {
      "name": "XX XX",
      "email": "XX.XX@XX.XX",
      "phone": "00",
      "fax": "00"
    }
  ]
}
```

### `www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/mensas/{id}/{language}/menus/daily/{date}/{daytime}`

Returns a list of meals available at a specific mensa at a specific date.

```
www.webservices.ethz.ch/gastro/v1/RVRI/Q1E1/mensas/12/de/menus/daily/2021-12-20/dinner
```
```json
{
  "id": 12,
  "mensa": "Mensa Polyterrasse",
  "daytime": "dinner",
  "hours": {
    "opening": [
      {
        "from": "11:15",
        "to": "14:15",
        "type": "opening_lunch"
      }
    ],
    "mealtime": [
      {
        "from": "11:15",
        "to": "14:15",
        "type": "lunch"
      }
    ]
  },
  "menu": {
    "date": "2021-12-20",
    "day": "Montag",
    "meals": [
      {
        "id": 90269,
        "mealtypes": [
          {
            "mealtype_id": 4,
            "label": "Fleisch"
          }
        ],
        "label": "HOME ABEND",
        "description": [
          "Mexican Meatballs",
          "Rindshackbällchen, Hot Salsa, ",
          "Couscous mit Mais und Peperonigemüse"
        ],
        "position": 0,
        "prices": {
          "student": "7.50",
          "staff": "9.90",
          "extern": "15.50"
        },
        "allergens": [
          {
            "allergen_id": 1,
            "label": "Gluten"
          }
        ],
        "origins": [
          {
            "origin_id": 39,
            "label": "Schweiz"
          }
        ]
      }
    ]
  }
}
```

## Glyph

The glyph API is used by the ETH App.

The Accept header has to be set to `application/json` to ensure that the response is JSON encoded (instead of XML).

Example:
```shell
$ curl http://glyph.ethz.ch/eth-ws/mensas -H "Accept: application/json"
```

### `glyph.ethz.ch/eth-ws/mensas`

Returns a list of all mensas

```
glyph.ethz.ch/eth-ws/mensas
```
```json
[
  {
    "mensaId": 28,
    "name": "food&lab",
    "location": "Zentrum",
    "openingTimes": "11:15–13:30",
    "imageUrl": "https://glyph.ethz.ch/eth-ws/resources/ETH_foodlab_03.jpg",
    "address": "CAB H 47.3\nUniversitätsstrasse 6\n8092 Zürich",
    "defaultMensa": true,
    "cams": []
  }
]
```

### `glyph.ethz.ch/eth-ws/mensas/detail/{id}`

Returns the list of currently available meals at the specified mensa.

```
glyph.ethz.ch/eth-ws/mensas/detail/12
```
```json
{
  "mensaId": 12,
  "name": "Mensa Polyterrasse",
  "location": "Zentrum",
  "openingTimes": "11:15–13:30",
  "mealTimes": "11:15–13:30",
  "web": "https://www.ethz.ch/de/campus/erleben/gastronomie-und-einkaufen/gastronomie/restaurants-und-cafeterias/zentrum/mensa-polyterrasse.html",
  "feedbackUrl": "https://www.ethz.ch/de/campus/erleben/gastronomie-und-einkaufen/gastronomie/gaeste-feedback.html",
  "meals": [
    {
      "type": "lunch",
      "menus": [
        {
          "menuId": 90064,
          "title": "LOCAL",
          "description": "Dieses Menu servieren wir Ihnen gerne bald wieder!\n",
          "price": {
            "student": "0.00",
            "staff": "2.00",
            "extern": "6.00"
          },
          "mealTypeId": 1,
          "mealType": "Sonstiges",
          "highlightAllergene": false
        }
      ]
    }
  ],
  "imageUrl": "https://glyph.ethz.ch/eth-ws/resources/ETH_mensa-polyterasse_06.jpg",
  "address": "MM B 79\nLeonhardstrasse 34\n8092 Zürich",
  "defaultMensa": true,
  "cams": []
}
```