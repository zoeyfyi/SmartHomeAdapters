package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.github.halspals.smarthomeadapters.smarthomeadapters.R
import com.google.gson.annotations.SerializedName

data class RobotStatus(
        @SerializedName("value")
        val value: Boolean,
        @SerializedName("current")
        val current: Int,
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
        const val TYPE_TOGGLE = "toggle"
        const val TYPE_RANGE = "range"
    }
}

data class RobotRegistrationBody(val nickname: String, val robotType: String)
