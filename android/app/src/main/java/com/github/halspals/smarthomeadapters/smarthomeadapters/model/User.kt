package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import com.google.gson.annotations.SerializedName

data class User(
    @SerializedName("name")
    val realName: String
)