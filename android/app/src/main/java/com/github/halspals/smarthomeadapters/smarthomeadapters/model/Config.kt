package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

/**
 * Model for the fine details of a configuration parameter's possible values.
 *
 * @property default the default value of the configuration parameter
 * @property min the minimum value that the parameter can take on
 * @property max the maximum value that the parameter can take on
 */
data class ConfigDetails(
    @SerializedName("default")
    val default: Int,
    @SerializedName("min")
    val min: Int,
    @SerializedName("max")
    val max: Int
)

/**
 * Model for a configuration parameter for a robot.
 *
 * @property id the unique id of the configuration parameter
 * @property name the parameter's name
 * @property description a description of the parameter
 * @property type the parameter's expected response type
 * @property details fine details of the parameter's possible values
 */
data class ConfigParameter(
    @SerializedName("id")
    val id: String,
    @SerializedName("name")
    val name: String,
    @SerializedName("description")
    val description: String,
    @SerializedName("type")
    val type: String,
    @SerializedName("details")
    val details: ConfigDetails
) {
    companion object {
        const val BOOL_TYPE = "bool"
        const val INT_TYPE = "int"
    }
}

/**
 * Model for the the user's chosen setting for a configuration parameter.
 *
 * @property id the unique id of the configuration parameter
 * @property the value that the user has chosen for the parameter
 */
data class ConfigResult(
    @SerializedName("id")
    val id: String,
    @SerializedName("value")
    val value: String
)