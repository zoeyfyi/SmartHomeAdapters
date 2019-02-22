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