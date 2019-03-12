package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

data class UseCase(
    @SerializedName("id")
    val id: String,
    @SerializedName("name")
    val name: String,
    @SerializedName("parameters")
    val parameters: List<ConfigParameter>
)