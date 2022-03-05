// hasSkillCount returns a count of items in the passed in group that have a non null skill property.
// this indicates that the skill has been at least injected
function hasSkillCount(group) {
    var out = 0

    group.skills.forEach(t => {
        if (t.skill) {
            out++
        }
    })

    return out
}

// same has hasSkillCount except here we check to ensure that the skill has been trained to level 5
function levelVSkillCount(group) {
    var out = 0

    group.skills.forEach(t => {
        if (t.skill && t.skill.trained_skill_level == 5) {
            out++
        }
    })

    return out
}

// same has levelVSkillCount except we are returning the sum of all of the skills skillpoints_in_skill attributes
function levelVTotalSP(group) {
    var out = 0

    group.skills.forEach(t => {
        if (t.skill && t.skill.trained_skill_level == 5) {
            out+=t.skill.skillpoints_in_skill
        }
    })

    return out
}

// potentialSP is the product of the rank of the skill multiplied by 256000
function potentialSP(rank) {
    return rank * 256000
}