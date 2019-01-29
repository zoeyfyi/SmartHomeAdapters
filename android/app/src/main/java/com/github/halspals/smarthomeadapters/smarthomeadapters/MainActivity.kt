package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.design.widget.BottomNavigationView
import android.support.v4.app.Fragment
import android.util.Log
import android.view.MenuItem
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity :
        AppCompatActivity(),
        BottomNavigationView.OnNavigationItemSelectedListener
{

    private val tag = "MainActivity"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        bottom_nav_bar.setOnNavigationItemSelectedListener(this)

        // Start the robots fragment by default
        startFragment(RobotsFragment())
    }

    override fun onNavigationItemSelected(item: MenuItem): Boolean {
        when (item.itemId) {
            R.id.robots_nav -> startFragment(RobotsFragment())
            R.id.triggers_nav -> startFragment(TriggersFragment())
            R.id.settings_nav -> startFragment(SettingsFragment())
            else -> {
                Log.e(tag, "[onNavigationItemSelected] id ${item.itemId} not recognized.")
                return false
            }
        }
        return true
    }

    /**
     * Replaces the currently active fragment, if there is any to replace.
     *
     * @param fragment the Fragment to replace the currently active one with.
     */
    private fun startFragment(fragment: Fragment) {
        Log.d(tag, "[startFragment] Invoked")
        val fManager = supportFragmentManager
        val fTransaction = fManager.beginTransaction()
        fTransaction.replace(R.id.fragmentContainer, fragment)
        fTransaction.commit()
        Log.d(tag, "[startFragment] Committed transaction to fragment")
    }
}
