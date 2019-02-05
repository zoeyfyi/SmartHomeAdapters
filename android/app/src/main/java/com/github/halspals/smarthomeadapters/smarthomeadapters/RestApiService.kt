package com.github.halspals.smarthomeadapters.smarthomeadapters

import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Token
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import retrofit2.Call
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.Body
import retrofit2.http.POST

interface RestApiService {

    @POST("register")
    fun registerUser(@Body user: User): Call<User>

    @POST("login")
    fun loginUser(@Body user: User): Call<Token>

    companion object {
        fun new(): RestApiService {
            return Retrofit.Builder()
                    .addConverterFactory(GsonConverterFactory.create())
                    .baseUrl("http://client.api.halspals.co.uk/")
                    .build()
                    .create(RestApiService::class.java)
        }
    }
}