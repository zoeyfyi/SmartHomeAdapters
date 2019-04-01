package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import kotlinx.android.synthetic.main.fragment_nickname.*
import kotlinx.android.synthetic.main.activity_register_robot.*
import okhttp3.ResponseBody
import org.jetbrains.anko.clearTask
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.intentFor
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

const val RENAME_FLAG = "RenameRobot"

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

        if (arguments?.getBoolean(RENAME_FLAG, false) == true) {
            // If we have been instructed to handle the renaming of an existing robot rather
            // than the naming of a new one, set up the appropriate button state
            continue_button.isEnabled = false
            continue_button.visibility = View.GONE
            update_name_button.visibility = View.VISIBLE
            update_name_button.isEnabled = true
        }

        continue_button.setOnClickListener { _ ->
            parent.robotNickname = device_name_text_view.text.toString()
            parent.startFragment(SelectAttachFragment())
        }

        cancel_button.setOnClickListener { _ ->
            parent.startActivity(parent.intentFor<MainActivity>().clearTask())
        }

        update_name_button.setOnClickListener { _ ->
            renameRobot(device_name_text_view.text.toString())
        }

    }

    /**
     * Takes a name for a robot already registered to the user and asks the server to rename it.
     *
     * @param newName the robot's new name to set
     */
    private fun renameRobot(newName: String) {
        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "got null access token, ex: $ex")
            } else {

                progressBar.visibility = View.VISIBLE
                cancel_button.isEnabled = false
                update_name_button.isEnabled = false

                parent.restApiService
                        .renameRobot(parent.robotId, accessToken, mapOf("nickname" to newName))
                        .enqueue(object : Callback<ResponseBody> {

                    override fun onResponse(call: Call<ResponseBody>, response: Response<ResponseBody>) {
                        progressBar.visibility = View.GONE
                        cancel_button.isEnabled = true
                        update_name_button.isEnabled = true

                        if (response.isSuccessful) {
                            Log.v(fTag, "[renameRobot] Success")
                            parent.toast("Updated nickname")
                            parent.finish()
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "[renameRobot] got unsuccessful "
                                    + "response, error: $error")
                            if (error != null) {
                                parent.snackbar_layout.snackbar(error)
                            }
                        }
                    }

                    override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                        progressBar.visibility = View.GONE
                        cancel_button.isEnabled = true
                        update_name_button.isEnabled = true

                        val error = t.message
                        Log.e(fTag, "[renameRobot] FAILED, error: $error")
                        if (error != null) {
                            parent.snackbar_layout.snackbar(error)
                        }
                    }
                })
            }
        }
    }
}
