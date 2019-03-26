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
 * @property nickname name of the robot
 * @property iconDrawable icon that represents the robot
 */
data class Robot(
    val id: String,
    val nickname: String,
    val robotType: String,
    @SerializedName("interfaceType") val robotInterfaceType: String,
    @SerializedName("status") val robotStatus: RobotStatus,
    val iconDrawable: Int = R.drawable.basic_home // default icon
) {
    companion object {
        const val INTERFACE_TYPE_TOGGLE = "toggle"
        const val INTERFACE_TYPE_RANGE = "range"
        const val ROBOT_TYPE_SWITCH = "switch"
        const val ROBOT_TYPE_THERMOSTAT = "thermostat"
        val ADD_ROBOT = Robot("","","","",RobotStatus(false,0,0,0))
    }
}

data class RobotRegistrationBody(val nickname: String, val robotType: String)
