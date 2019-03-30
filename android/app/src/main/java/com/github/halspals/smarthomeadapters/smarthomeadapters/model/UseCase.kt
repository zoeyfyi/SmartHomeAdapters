package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

/**
 * Model for the use case / attachment of a robot
 *
 * @property id the identifier of the use case
 * @property name the use case's common name
 * @property parameters the [ConfigParameter]s associated with the use case
 */
data class UseCase(
    @SerializedName("id")
    val id: String,
    @SerializedName("name")
    val name: String,
    @SerializedName("parameters")
    val parameters: List<ConfigParameter>
)