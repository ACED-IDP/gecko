package fixtures

var TestConfig string = `[
    {
        "tabTitle": "test",
        "guppyConfig": {
            "dataType": "file",
            "nodeCountTitle": "file Count",
            "fieldMapping": []
        },
        "charts": {
            "a": {
                "chartType": "bar",
                "title": "a"
            },
            "b": {
                "chartType": "bar",
                "title": "a"
            }
        },
        "filters": {
            "tabs": [
                {
                    "title": "Filters",
                    "fields": [
                        "a",
                        "b",
                        "project_id"
                    ],
                    "fieldsConfig": {
                        "a": {
                            "field": "a",
                            "dataField": "",
                            "index": "",
                            "label": "a",
                            "type": "enum"
                        },
                        "b": {
                            "field": "b",
                            "dataField": "",
                            "index": "",
                            "label": "b",
                            "type": "enum"
                        },
                        "project_id": {
                            "field": "project_id",
                            "dataField": "",
                            "index": "",
                            "label": "Project ID",
                            "type": "enum"
                        }
                    }
                }
            ]
        },
        "table": {
            "enabled": true,
            "fields": [
                "project_id",
                "b",
                "a"
            ],
            "columns": {
                "project_id": {
                    "field": "project_id",
                    "title": "Project ID"
                },
                "b": {
                    "field": "b",
                    "title": "asd"
                },
                "a": {
                    "field": "a",
                    "title": "a"
                }
            }
        },
        "dropdowns": {},
        "buttons": [],
        "loginForDownload": false
    }
]`
