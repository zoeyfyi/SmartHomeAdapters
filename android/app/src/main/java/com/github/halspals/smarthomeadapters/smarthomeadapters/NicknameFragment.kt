package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_nickname.*
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.intentFor

/**
 * A fragment which allows the user to set the nickname for a robot being registered or edited.
 */
class NicknameFragment : Fragment() {

    private val fTag = "NicknameFragment"

    private val parent by lazy { activity as RegisterRobotActivity }

    override fun onCreateView(
            inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_nickname, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)


        continue_button.setOnClickListener { _ ->
            parent.robotNickname = device_name_text_view.text.toString()
            parent.startFragment(SelectAttachFragment())
        }

        cancel_button.setOnClickListener { _ ->
            parent.startActivity(parent.intentFor<MainActivity>().clearTask())
        }

    }
}
