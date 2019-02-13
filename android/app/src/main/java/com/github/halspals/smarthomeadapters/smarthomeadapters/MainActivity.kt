package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.os.Bundle
import android.support.design.widget.BottomNavigationView
import android.support.v4.app.Fragment
import android.support.v4.app.FragmentManager
import android.util.Log
import android.view.MenuItem
import com.google.zxing.integration.android.IntentIntegrator
import com.google.zxing.integration.android.IntentResult
import kotlinx.android.synthetic.main.activity_main.*

class MainActivity :
        AppCompatActivity(),
        BottomNavigationView.OnNavigationItemSelectedListener
{

    internal val restApiService by lazy {
        RestApiService.new()
    }

    private val tag = "MainActivity"

    private var authenticationToken: String? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        bottom_nav_bar.setOnNavigationItemSelectedListener(this)

        // Start the robots fragment by default
        startFragment(RobotsFragment())

        // Get the token for the current session
        authenticationToken = intent.getStringExtra("token")
        Log.d(tag, "token: $authenticationToken")
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
     * @param addToBackstack if true, fragment will be added to the backstack,
     * otherwise backstack will be dropped
     */
    fun startFragment(fragment: Fragment, addToBackstack: Boolean = false) {
        Log.d(tag, "[startFragment] Invoked")

        val fManager = supportFragmentManager
        fManager.beginTransaction().run {
            replace(R.id.fragmentContainer, fragment)

            // manually handle the backstack
            if (addToBackstack) {
                // A->B to A->B->C (add to backstack)
                addToBackStack(null)
            } else {
                // A->B->C to A (clear backstack)
                fManager.popBackStack(null, FragmentManager.POP_BACK_STACK_INCLUSIVE)
            }

            commit()
        }

        Log.d(tag, "[startFragment] Committed transaction to fragment")
    }
}
