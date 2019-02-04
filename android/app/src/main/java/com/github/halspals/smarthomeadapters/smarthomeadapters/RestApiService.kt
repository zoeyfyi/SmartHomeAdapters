package com.github.halspals.smarthomeadapters.smarthomeadapters

import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Token
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import io.reactivex.Observable
import retrofit2.Retrofit
import retrofit2.adapter.rxjava2.RxJava2CallAdapterFactory
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.Body
import retrofit2.http.POST

interface RestApiService {

    @POST("register")
    fun registerUser(@Body user: User): Observable<User>

    @POST("login")
    fun loginUser(@Body user: User): Observable<Token>

    companion object {
        fun create(): RestApiService {
            return Retrofit.Builder()
                    .addCallAdapterFactory(RxJava2CallAdapterFactory.create())
                    .addConverterFactory(GsonConverterFactory.create())
                    .baseUrl("http://client.api.halspals.co.uk/")
                    .build()
                    .create(RestApiService::class.java)
        }
    }
}