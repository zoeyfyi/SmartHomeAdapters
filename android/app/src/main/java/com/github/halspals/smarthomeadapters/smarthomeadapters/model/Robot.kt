package com.github.halspals.smarthomeadapters.smarthomeadapters.model

sealed class RobotInterface {
    class Toggle(val isOn: Boolean) : RobotInterface()
    class Range(val value: Int, val min: Int, val max: Int) : RobotInterface()
}

/**
 * Model for a smart home adapter robot
 *
 * @property nickname name of the robot
 * @property iconDrawable icon that represents the robot
 */
data class Robot(
    val id: String,
    val nickname: String,
    val iconDrawable: Int,
    val robotInterface: RobotInterface
)