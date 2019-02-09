package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_settings.*
import org.jetbrains.anko.toast

class SettingsFragment : Fragment() {

    private val fTag = "RobotFragment"

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_settings, container, false)

    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        sign_out_button.setOnClickListener { _ -> deletePrefsAndSignOut(view.context) }
    }

    /**
     * Clears the shared preferences file and invokes [MainActivity.moveToAuthenticationActivity].
     */
    private fun deletePrefsAndSignOut(context: Context) {
        val prefs = context.getSharedPreferences(SHARED_PREFERENCES_FILE, Context.MODE_PRIVATE)
        val editor = prefs.edit()
        editor.clear()  // Delete all the preferences stored in the file
        editor.apply()

        context.toast("Signed out")

        (activity as MainActivity).moveToAuthenticationActivity()
    }
}
