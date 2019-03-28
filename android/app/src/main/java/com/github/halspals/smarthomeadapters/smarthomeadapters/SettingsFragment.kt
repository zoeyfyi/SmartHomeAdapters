package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_settings.*
import net.openid.appauth.AuthState
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor

class SettingsFragment : Fragment() {

    private val fTag = "RobotFragment"

    private val parent by lazy { activity as MainActivity }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_settings, container, false)

    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        logout_button.setOnClickListener { _ ->
            // The user has asked to be signed out; write a clean AuthState to the device
            // and start the authentication activity
            writeAuthState(parent, AuthState())
            parent.startActivity(parent.intentFor<AuthenticationActivity>().clearTask())
        }
    }
}
