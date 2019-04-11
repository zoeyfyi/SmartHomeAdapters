package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

/**
 * Model for the use case / attachment of a robot
 *
 * @property name the use case's common name
 * @property parameters the [ConfigParameter]s associated with the use case
 */
data class UseCase(
    @SerializedName("name")
    val name: String,
    @SerializedName("description")
    val description: String
)