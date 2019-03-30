package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.github.halspals.smarthomeadapters.smarthomeadapters.R
import com.google.gson.annotations.SerializedName

data class RobotStatus(
        @SerializedName("value")
        var value: Boolean,
        @SerializedName("current")
        var current: Int,
        @SerializedName("min")
        val min: Int,
        @SerializedName("max")
        val max: Int
)

/**
 * Model for a smart home adapter robot
 *
 * @property id the robot's unique ID
 * @property nickname name of the robot
 * @property robotType the type of functionality the robot provides (toggle, range, etc)
 * @property robotInterfaceType the type of interface/use case of the robot (switch, thermostat, etc)
 * @property robotStatus the robot's current state
 */
data class Robot(
    val id: String,
    val nickname: String,
    val robotType: String,
    @SerializedName("interfaceType") val robotInterfaceType: String,
    @SerializedName("status") val robotStatus: RobotStatus
) {
    companion object {
        const val INTERFACE_TYPE_TOGGLE = "toggle"
        const val INTERFACE_TYPE_RANGE = "range"
        const val ROBOT_TYPE_SWITCH = "switch"
        const val ROBOT_TYPE_THERMOSTAT = "thermostat"
        val ADD_ROBOT = Robot("","","","",
                RobotStatus(false,0,0,0))
    }
}

/**
 * Model for the body of a robot registration call
 *
 * @property nickname the nickname to set for the robot
 * @property robotType the type of functionality the robot provides
 */
data class RobotRegistrationBody(val nickname: String, val robotType: String)
