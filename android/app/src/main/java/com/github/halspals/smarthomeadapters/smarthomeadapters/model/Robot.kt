package com.github.halspals.smarthomeadapters.smarthomeadapters.model

/**
 * Model for a smart home adapter robot
 *
 * @property nickname name of the robot
 * @property iconDrawable icon that represents the robot
 */
data class Robot(
    val nickname: String,
    val iconDrawable: Int
)