

import logging
import time
import json
import uuid
import requests

# Imports for v3 validation
from validation import validate_message

# Setup logger
logger = logging.getLogger()
logger.setLevel(logging.INFO)


SAMPLE_APPLIANCES=[]

def lambda_handler(request, context):
    """Main Lambda handler.

    Since you can expect both v2 and v3 directives for a period of time during the migration
    and transition of your existing users, this main Lambda handler must be modified to support
    both v2 and v3 requests.
    """

    try:
        logger.info("Directive:")
        logger.info(json.dumps(request, indent=4, sort_keys=True))



        logger.info("Received v3 directive!")
        if request["directive"]["header"]["name"] == "Discover":
            response = handle_discovery_v3(request)
        else:
            response = handle_non_discovery_v3(request)

        logger.info("Response:")
        logger.info(json.dumps(response, indent=4, sort_keys=True))

        #if version == "3":
            #logger.info("Validate v3 response")
            #validate_message(request, response)

        return response
    except ValueError as error:
        logger.error(error)
        raise


# v2 utility functions
def get_appliance_by_appliance_id(appliance_id):
    for appliance in SAMPLE_APPLIANCES:
        if appliance["applianceId"] == appliance_id:
            return appliance
    return None

def get_utc_timestamp(seconds=None):
    return time.strftime("%Y-%m-%dT%H:%M:%S.00Z", time.gmtime(seconds))

def get_uuid():
    return str(uuid.uuid4())

# v3 handlers
def handle_discovery_v3(request):
    endpoints = []
    token = request["directive"]["payload"]["scope"]["token"]
    repon = requests.get("https://client.api.test.halspals.co.uk/robots", headers = {"token":token})
    getdevice = json.loads(repon.text)


    for i in getdevice:
        if i["robotType"] == "switch":
            SAMPLE_APPLIANCES.append(
        {
        "applianceId": i["id"],
        "manufacturerName": "Sample Manufacturer",
        "modelName": "Smart Switch",
        "version": "1",
        "friendlyName": i["nickname"],
        "friendlyDescription": "Switch that can only be turned on/off",
        "isReachable": True,
        "actions": [
            "turnOn",
            "turnOff"
        ],
        "additionalApplianceDetails": {}
        })
        elif i["robotType"] == "thermostat":
            SAMPLE_APPLIANCES.append(
        {
        "applianceId": i["id"],
        "manufacturerName": "Sample Manufacturer",
        "modelName": "Smart Thermostat",
        "version": "1",
        "friendlyName": i["nickname"],
        "friendlyDescription": "Thermostat that can change and query temperatures of thermostat",
        "isReachable": True,
        "actions": [
            "setTargetTemperature",
            "incrementTargetTemperature",
            "decrementTargetTemperature",
            "getTargetTemperature"
        ],
        "additionalApplianceDetails": {}
        })
        elif i["robotType"] == "boltlock":
            SAMPLE_APPLIANCES.append(
        {
        "applianceId": i["id"],
        "manufacturerName": "Sample Manufacturer",
        "modelName": "Smart Switch",
        "version": "1",
        "friendlyName": i["nickname"],
        "friendlyDescription": "Lock that can be locked and can query lock state",
        "isReachable": True,
        "actions": [
            "turnOn",
            "turnOff"
        ],
        "additionalApplianceDetails": {}
        })


    for appliance in SAMPLE_APPLIANCES:
        endpoints.append(get_endpoint_from_v2_appliance(appliance))


    response = {
        "event": {
            "header": {
                "namespace": "Alexa.Discovery",
                "name": "Discover.Response",
                "payloadVersion": "3",
                "messageId": get_uuid()
            },
            "payload": {
                "endpoints": endpoints
            }
        }
    }
    return response

def handle_non_discovery_v3(request):
    request_namespace = request["directive"]["header"]["namespace"]
    request_name = request["directive"]["header"]["name"]


    if request_namespace == "Alexa.PowerController":
        request_equ_id = request["directive"]["endpoint"]["endpointId"]
        token_needed = request["directive"]["endpoint"]["scope"]["token"]
        if request_name == "TurnOn":
            respones_test = requests.patch("https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/toggle/true", headers = {"token":token_needed})
            value = "ON"
        else:
            respones_test = requests.patch(
                "https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/toggle/false", headers={"token":token_needed})
            value = "OFF"

        response = {
            "context": {
                "properties": [
                    {
                        "namespace": "Alexa.PowerController",
                        "name": "powerState",
                        "value": value,
                        "timeOfSample": get_utc_timestamp(),
                        "uncertaintyInMilliseconds": 500
                    }
                ]
            },
            "event": {
                "header": {
                    "namespace": "Alexa",
                    "name": "Response",
                    "payloadVersion": "3",
                    "messageId": get_uuid(),
                    "correlationToken": request["directive"]["header"]["correlationToken"]
                },
                "endpoint": {
                    "scope": {
                        "type": "BearerToken",
                        "token": token_needed
                    },
                    "endpointId": request_equ_id
                },
                "payload": {}
            }
        }
        return response
    elif request_namespace == "Alexa":
        request_equ_id = request["directive"]["endpoint"]["endpointId"]
        token_needed = request["directive"]["endpoint"]["scope"]["token"]

        get_current = requests.get("https://client.api.test.halspals.co.uk/robot/" + request_equ_id,
                                   headers={"token": token_needed})
        get_current = json.loads(get_current.text)

        get_current_temp = get_current["status"]["current"]



        if request_name == "ReportState":

            response = {
                    "context": {
                    "properties": [ {
                    "namespace": "Alexa.TemperatureSensor",
                    "name": "temperature",
                     "value": {
                          "value": get_current_temp,
                          "scale": "KELVIN"
                           },
                     "timeOfSample": get_utc_timestamp(),
                    "uncertaintyInMilliseconds": 1000
                    },{
                "namespace": "Alexa.ThermostatController",
                "name": "thermostatMode",
                "value": "HEAT",
                "timeOfSample": get_utc_timestamp(),
                "uncertaintyInMilliseconds": 6000
                  }, ]
                  },
                 "event": {
                "header": {
               "namespace": "Alexa",
               "name": "StateReport",
                 "payloadVersion": "3",
                  "messageId": get_uuid(),
                 "correlationToken": request["directive"]["header"]["correlationToken"]
                    },
                  "endpoint": {
                  "endpointId": request_equ_id,
                      "scope": {
                          "type": "BearerToken",
                          "token": "access-token-from-Amazon"
                      },
                      "cookie": {}
                  },
                "payload": {}
                  }
            }
            return response



    elif request_namespace == "Alexa.LockController":
        request_equ_id = request["directive"]["endpoint"]["endpointId"]
        token_needed = request["directive"]["endpoint"]["scope"]["token"]

        if request_name == "Lock":
            respones_test = requests.patch(
                "https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/toggle/true",
                headers={"token": token_needed})
            value = "LOCKED"
        elif request_name == "Unlock":
            respones_test = requests.patch(
                "https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/toggle/false",
                headers={"token": token_needed})
            value = "UNLOCKED"

        response = {
            "event": {
                "header": {
                    "namespace": "Alexa",
                    "name": "Response",
                    "payloadVersion": "3",
                    "messageId": get_uuid(),
                    "correlationToken": request["directive"]["header"]["correlationToken"]
                },
                "endpoint": {
                    "scope": {
                        "type": "BearerToken",
                        "token": token_needed
                    },
                    "endpointId":request_equ_id
                },
                "payload": {}
            },
            "context": {
                "properties": [
                    {
                        "namespace": "Alexa.LockController",
                        "name": "lockState",
                        "value": value,
                        "timeOfSample": get_utc_timestamp(),
                        "uncertaintyInMilliseconds": 1000
                    }
                ]
            },
        }
        return response

    elif request_namespace == "Alexa.ThermostatController":
        request_equ_id = request["directive"]["endpoint"]["endpointId"]
        token_needed = request["directive"]["endpoint"]["scope"]["token"]


        if request_name == "SetTargetTemperature":
            targettemp = int(request["directive"]["payload"]["targetSetpoint"]["value"])
            targetscale = request["directive"]["payload"]["targetSetpoint"]["scale"]

            if targetscale == "CELSIUS":
                targettemp = str(targettemp + 273)
            elif targetscale == "FAHRENHEIT":
                targettemp = str(targettemp + 255)
            else:
                targettemp = str(targettemp)

            respones_test = requests.patch(
                "https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/range/" + targettemp,
                headers={"token": token_needed})

            response = {
                "context": {
                    "properties": [{
                        "namespace": "Alexa.ThermostatController",
                        "name": "targetSetpoint",
                        "value": {
                            "value": targettemp,
                            "scale": "KELVIN"
                        },
                        "timeOfSample": get_utc_timestamp(),
                        "uncertaintyInMilliseconds": 500
                    }]
                },
                "event": {
                    "header": {
                        "namespace": "Alexa",
                        "name": "Response",
                        "payloadVersion": "3",
                        "messageId": get_uuid(),
                        "correlationToken": request["directive"]["header"]["correlationToken"]
                    },
                    "endpoint": {
                        "endpointId": request_equ_id
                    },
                    "payload": {}
                }
            }
            return response

        elif request_name == "AdjustTargetTemperature":
            targettemp = int(request["directive"]["payload"]["targetSetpointDelta"]["value"])
            targetscale = request["directive"]["payload"]["targetSetpointDelta"]["scale"]

            get_current = requests.get("https://client.api.test.halspals.co.uk/robot/" + request_equ_id, headers = {"token":token_needed})
            logger.info(json.dumps(get_current.text, indent=4, sort_keys=True))
            get_current = json.loads(get_current.text)

            get_current_temp = get_current["status"]["current"]

            if targetscale == "CELSIUS":
                targettemp = str(targettemp + get_current_temp)
            elif targetscale == "FAHRENHEIT":
                targettemp = str(targettemp + get_current_temp)
            else:
                targettemp = str(targettemp + get_current_temp)



            respones_test = requests.patch(
                "https://client.api.test.halspals.co.uk/robot/" + request_equ_id + "/range/" + targettemp,
                headers={"token": token_needed})

            response = {
                "context": {
                    "properties": {
                        "namespace": "Alexa.ThermostatController",
                        "name": "targetSetpoint",
                        "value": {
                            "value": targettemp,
                            "scale": "KELVIN"
                        },
                        "timeOfSample": get_utc_timestamp(),
                        "uncertaintyInMilliseconds": 500
                    }
                },
                "event": {
                    "header": {
                        "namespace": "Alexa",
                        "name": "Response",
                        "payloadVersion": "3",
                        "messageId": get_uuid(),
                        "correlationToken": request["directive"]["header"]["correlationToken"]
                    },
                    "endpoint": {
                        "endpointId": request_equ_id
                    },
                    "payload": {}
                }
            }
            return response

        elif request_name == "SetThermostatMode":
            moderequre = request["directive"]["payload"]["thermostatMode"]["value"]

            response = {
                "context": {
                     "properties":  {
                          "namespace": "Alexa.ThermostatController",
                          "name": "thermostatMode",
                          "value": "HEAT",
                           "timeOfSample": get_utc_timestamp(),
                          "uncertaintyInMilliseconds": 500
                          }
                          },
                     "event": {
                          "header": {
                              "namespace": "Alexa",
                              "name": "Response",
                               "payloadVersion": "3",
                                "messageId": get_uuid(),
                               "correlationToken": request["directive"]["header"]["correlationToken"]
                          },
                            "endpoint": {
                           "endpointId": request_equ_id
                           },
                           "payload": {}
                      }
            }
            return response



    elif request_namespace == "Alexa.Authorization":
        if request_name == "AcceptGrant":
            response = {
                "event": {
                    "header": {
                        "namespace": "Alexa.Authorization",
                        "name": "AcceptGrant.Response",
                        "payloadVersion": "3",
                        "messageId": get_uuid()
                    },
                    "payload": {}
                }
            }
            return response

    # other handlers omitted in this example

# v3 utility functions
def get_endpoint_from_v2_appliance(appliance):
    endpoint = {
        "endpointId": appliance["applianceId"],
        "manufacturerName": appliance["manufacturerName"],
        "friendlyName": appliance["friendlyName"],
        "description": appliance["friendlyDescription"],
        "displayCategories": [],
        "cookie": appliance["additionalApplianceDetails"],
        "capabilities": []
    }
    endpoint["displayCategories"] = get_display_categories_from_v2_appliance(appliance)
    endpoint["capabilities"] = get_capabilities_from_v2_appliance(appliance)
    return endpoint

def get_directive_version(request):
    try:
        return request["directive"]["header"]["payloadVersion"]
    except:
        try:
            return request["header"]["payloadVersion"]
        except:
            return "-1"

def get_endpoint_by_endpoint_id(endpoint_id):
    appliance = get_appliance_by_appliance_id(endpoint_id)
    if appliance:
        return get_endpoint_from_v2_appliance(appliance)
    return None

def get_display_categories_from_v2_appliance(appliance):
    model_name = appliance["modelName"]
    if model_name == "Smart Switch": displayCategories = ["SWITCH"]
    elif model_name == "Smart Light": displayCategories = ["LIGHT"]
    elif model_name == "Smart White Light": displayCategories = ["LIGHT"]
    elif model_name == "Smart Thermostat": displayCategories = ["THERMOSTAT"]
    elif model_name == "Smart Lock": displayCategories = ["SMARTLOCK"]
    elif model_name == "Smart Scene": displayCategories = ["SCENE_TRIGGER"]
    elif model_name == "Smart Activity": displayCategories = ["ACTIVITY_TRIGGER"]
    elif model_name == "Smart Camera": displayCategories = ["CAMERA"]
    else: displayCategories = ["OTHER"]
    return displayCategories

def get_capabilities_from_v2_appliance(appliance):
    model_name = appliance["modelName"]
    if model_name == 'Smart Switch':
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerState" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            }
        ]
    elif model_name == "Smart Light":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerState" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.ColorController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "color" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.ColorTemperatureController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "colorTemperatureInKelvin" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.BrightnessController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "brightness" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerLevelController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerLevel" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PercentageController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "percentage" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            }
        ]
    elif model_name == "Smart White Light":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerState" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.ColorTemperatureController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "colorTemperatureInKelvin" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.BrightnessController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "brightness" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerLevelController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerLevel" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PercentageController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "percentage" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            }
        ]
    elif model_name == "Smart Thermostat":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.ThermostatController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "targetSetpoint" },
                        { "name": "thermostatMode" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.TemperatureSensor",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "temperature" }
                    ],
                    "proactivelyReported": False,
                    "retrievable": True
                }
            }

        ]
    elif model_name == "Smart Thermostat Dual":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.ThermostatController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "upperSetpoint" },
                        { "name": "lowerSetpoint" },
                        { "name": "thermostatMode" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                },
                "configuration": {
                    "supportsScheduling": False,
                    "supportedModes": [
                        "HEAT",
                    ]
                }
            },
            {
                "type": "AlexaInterface",
                "interface": "Alexa.TemperatureSensor",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "temperature" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            }
        ]
    elif model_name == "Smart Lock":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.LockController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "lockState" }
                    ],
                    "proactivelyReported": True,
                    "retrievable": True
                }
            }
        ]
    elif model_name == "Smart Scene":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.SceneController",
                "version": "3",
                "supportsDeactivation": False,
                "proactivelyReported": True
            }
        ]
    elif model_name == "Smart Activity":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.SceneController",
                "version": "3",
                "supportsDeactivation": True,
                "proactivelyReported": True
            }
        ]
    elif model_name == "Smart Camera":
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.CameraStreamController",
                "version": "3",
                "cameraStreamConfigurations" : [ {
                    "protocols": ["RTSP"],
                    "resolutions": [{"width":1280, "height":720}],
                    "authorizationTypes": ["NONE"],
                    "videoCodecs": ["H264"],
                    "audioCodecs": ["AAC"]
                } ]
            }
        ]
    else:
        # in this example, just return simple on/off capability
        capabilities = [
            {
                "type": "AlexaInterface",
                "interface": "Alexa.PowerController",
                "version": "3",
                "properties": {
                    "supported": [
                        { "name": "powerState" }
                    ],
                    "proactivelyReported": False,
                    "retrievable": True
                }
            }
        ]

    # additional capabilities that are required for each endpoint
    endpoint_health_capability = {
        "type": "AlexaInterface",
        "interface": "Alexa.EndpointHealth",
        "version": "3",
        "properties": {
            "supported":[
                { "name":"connectivity" }
            ],
            "proactivelyReported": True,
            "retrievable": True
        }
    }
    alexa_interface_capability = {
        "type": "AlexaInterface",
        "interface": "Alexa",
        "version": "3"
    }
    capabilities.append(endpoint_health_capability)
    capabilities.append(alexa_interface_capability)
    return capabilities
