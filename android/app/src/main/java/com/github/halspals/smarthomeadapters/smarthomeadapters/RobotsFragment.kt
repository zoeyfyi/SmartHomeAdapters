package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.GridView
import android.widget.PopupMenu
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.User
import kotlinx.android.synthetic.main.activity_main.*
import kotlinx.android.synthetic.main.fragment_robots.*
import org.jetbrains.anko.design.snackbar
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

/**
 * The main screen of the app, presenting the user with their robots and allowing interaction thereof.
 */
class RobotsFragment : Fragment() {

    private val fTag = "RobotsFragment"

    private lateinit var robotGrid: GridView

    private val parent by lazy { activity as MainActivity }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {

        Log.d(fTag, "[onCreateView] Invoked")
        return inflater.inflate(R.layout.fragment_robots, container, false)
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken, _, ex ->
            if (accessToken == null) {
                Log.e(fTag, "[performActionWithFreshTokens] got null access token, "
                        + "exception: $ex")
            } else {
                fetchRobots(accessToken, view)
                getUserName(accessToken)
            }
        }

        edit_mode_image_view.setOnClickListener { _ ->
            parent.isInEditMode = true
            edit_mode_image_view.visibility = View.INVISIBLE
            more_options_image_view.visibility = View.INVISIBLE
            finish_edit_image_view.visibility = View.VISIBLE
        }

        finish_edit_image_view.setOnClickListener { _ ->
            parent.isInEditMode = false
            edit_mode_image_view.visibility = View.VISIBLE
            more_options_image_view.visibility = View.VISIBLE
            finish_edit_image_view.visibility = View.INVISIBLE
        }

        more_options_image_view.setOnClickListener { _ ->
            Log.v(fTag, "Inflating settings_menu")
            PopupMenu(more_options_image_view.context, more_options_image_view).run {
                menuInflater.inflate(R.menu.settings_menu, menu)
                setOnMenuItemClickListener { menuItem ->
                    when (menuItem.itemId) {
                        R.id.log_out_item -> {
                            Log.v(fTag, "Logging out the user")
                            parent.signOut()
                            true
                        }
                        else -> {
                            Log.e(fTag, "Unexpected menu click on item " +
                                    "with ID ${menuItem.itemId}")
                            false
                        }
                    }
                }
                show()
            }
        }
    }

    /**
     * Fetches all the robots associated with the user and invokes [displayRobots].
     */
    private fun fetchRobots(token: String, view: View) {
        parent.restApiService
                .getRobots(token)
                .enqueue(object : Callback<List<Robot>> {

                    override fun onFailure(call: Call<List<Robot>>, t: Throwable) {
                        val errorMsg = t.message
                        Log.e(fTag, "[getRobots] FAILED, got error: $errorMsg")
                        if (errorMsg != null) {
                            parent.snackbar_layout.snackbar(errorMsg)
                        }

                        displayRobots(view, listOf())
                    }

                    override fun onResponse(
                            call: Call<List<Robot>>,
                            response: Response<List<Robot>>) {

                        val robots = response.body()
                        val robotsToList = if (response.isSuccessful && robots != null) {
                            Log.v(fTag, "[getRobots] got robots: $robots")
                            robots
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "[getRobots] got unsuccessful response, error: $error")
                            if (error != null) {
                                parent.snackbar_layout.snackbar(error)
                            }

                            listOf()
                        }

                        displayRobots(view, robotsToList)
                    }
                })
    }

    /**
     * Sets up the [robotGrid] with the robots associated with the user.
     */
    private fun displayRobots(view: View, robots: List<Robot>) {
        robotGrid = view.findViewById(R.id.RobotGrid)
        robotGrid.adapter = RobotAdapter(parent, robots.toMutableList())
    }

    private fun getUserName(token: String) {
        parent.restApiService
                .getUserName(token)
                .enqueue(object : Callback<User> {
                    override fun onResponse(call: Call<User>, response: Response<User>) {
                        val user = response.body()

                        if (response.isSuccessful && user != null) {
                            Log.d(fTag, "[getUserName] got User w/ real name ${user.realName}")
                            welcome_text_view.text = getString(R.string.welcome_text, user.realName)
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "[getUserName] unsuccessful response or null user; " +
                                    "user was $user, error was $error")
                        }

                    }

                    override fun onFailure(call: Call<User>, t: Throwable) {
                        val errorMsg = t.message
                        Log.e(fTag, "[getUserName] FAILED, got error: $errorMsg")
                        if (errorMsg != null) {
                            parent.snackbar_layout.snackbar(errorMsg)
                        }
                    }

                })
    }

}
