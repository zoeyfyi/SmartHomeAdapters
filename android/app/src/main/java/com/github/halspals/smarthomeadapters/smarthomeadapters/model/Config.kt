package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

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
 * @property name the parameter's name
 * @property description a description of the parameter
 * @property type the parameter's expected response type
 * @property details fine details of the parameter's possible values
 */
data class ConfigParameter(
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

data class ConfigResult(
    @SerializedName("name")
    val name: String,
    @SerializedName("value")
    val value: String
)