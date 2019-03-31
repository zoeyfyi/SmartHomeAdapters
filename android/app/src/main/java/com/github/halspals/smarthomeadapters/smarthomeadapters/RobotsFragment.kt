package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.os.Bundle
import android.support.v4.app.Fragment
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.GridView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
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
                        Log.e(fTag, "getRobots FAILED, got error: $errorMsg")
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
                            robots
                        } else {
                            val error = RestApiService.extractErrorFromResponse(response)
                            Log.e(fTag, "getRobots got unsuccessful response, error: $error")
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


}
