package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.text.Editable
import android.text.TextWatcher
import android.widget.Button

internal class NonEmptyTextWatcher(
        private val buttonUpdater: ButtonUpdater,
        private val targetButton: Button)
    : TextWatcher {

    override fun afterTextChanged(s: Editable?) {
        // Not interested
    }

    override fun beforeTextChanged(s: CharSequence?, start: Int, count: Int, after: Int) {
        // Not interested
    }

    override fun onTextChanged(s: CharSequence?, start: Int, before: Int, count: Int) {
        buttonUpdater.updateButton(targetButton)
    }
}