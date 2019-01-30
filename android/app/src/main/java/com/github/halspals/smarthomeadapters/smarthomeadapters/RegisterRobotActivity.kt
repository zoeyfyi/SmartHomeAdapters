package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log

class RegisterRobotActivity : AppCompatActivity() {

    val tag = "RegisterRobotActivity "

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_register_robot)
        startFragment(ConnectToAdapterFragment())
    }

    fun startFragment(fragment: Fragment) {
        Log.d(tag, "[startFragment] Invoked")

        val fManager = supportFragmentManager
        fManager.beginTransaction().run {
            replace(R.id.fragment_container, fragment)
            commit()
        }

        Log.d(tag, "[startFragment] Committed transaction to fragment")
    }

}
