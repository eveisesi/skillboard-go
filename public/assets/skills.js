let selectedSkillGroupID = -1

// Max 5, Min -2, -2 == All Skills, -1 == All Uninjected Skiills, >= 0 == All Injected Skills
let selectedSkillLevel = 0

function setActiveGroupID(clickedGroupID) {
    selectedSkillGroupID = parseInt(clickedGroupID)
    renderSkillGroupListGroup()
}


function setFilterSkillLevel(skillLevel) {
    selectedSkillLevel = parseInt(skillLevel)
    renderSkillLevelListGroup()
    renderSkillGroupDetails()
}

function renderSkillLevelListGroup() {
    let rows = []
    rows.push(`
    <li data-level="-2" class="list-group-item text-white ${selectedSkillLevel === -2 ? 'active' : ''}" onclick="setFilterSkillLevel(this.getAttribute('data-level'))">
    All Skills
    </li>
    `, `
    <li data-level="-1" class="list-group-item text-white ${selectedSkillLevel === -1 ? 'active' : ''}" onclick="setFilterSkillLevel(this.getAttribute('data-level'))">
    All Uninjected/Untrained Skills
    </li>\
    `, `
    <li data-level="0" class="list-group-item text-white ${selectedSkillLevel === 0 ? 'active' : ''}" onclick="setFilterSkillLevel(this.getAttribute('data-level'))">
    All Injected Skills
    </li>
    `)
    for (let i = 1; i <= 5; i++) {
        rows.push(`
        <li data-level="${i}" class="list-group-item text-white ${selectedSkillLevel === i ? 'active' : ''}" onclick="setFilterSkillLevel(this.getAttribute('data-level'))">
        Skill Trained Level ${i}
        </li>
        `)
    }
    document.querySelector("#skillLevelListGroup").innerHTML = `<ul class="list-group">${rows.join("")}</ul>`;
}
function renderSkillGroupListGroup() {
    let rows = skillzGrouped.map(group => {
        return `
            <li data-id="${group.id}" class="list-group-item text-white d-flex w-100 justify-content-between ${group.id == selectedSkillGroupID ? 'active' : ''}" onclick="setActiveGroupID(this.getAttribute('data-id'))">
            <span>
            ${group.name}
            </span>
            <span>
                (${hasSkillCount(group)} / ${numeral(group.totalGroupSP).format("0,0")})
            </span>
            </li>
        `
    })
    rows.unshift(
        `<li data-id="-1" class="list-group-item text-white ${selectedSkillGroupID === -1 ? 'active' : ''}" onclick="setActiveGroupID(this.getAttribute('data-id'))">All</li>`
    )
    document.querySelector("#skillGroupListGroup").innerHTML = `<ul class="list-group">${rows.join("")}</ul>`;
    renderSkillGroupDetails()

}


function renderSkillGroupDetails() {
    const groups = skillzGrouped
    let cards = [];
    for (let i = 0; i < groups.length; i++) {
        const group = groups[i]
        if (selectedSkillGroupID > 0 && selectedSkillGroupID != group.id) {
            continue
        }
        const skillLines = []
        for (let i = 0; i < group.skills.length; i++) {
            const skill = group.skills[i]
            let typeName = `<span class="text-muted"><strike>${skill.name}</strike></span>`
            if (skill.skill) {
                typeName = `<span>${skill.name}</span>`
            }
            let rank = 0
            let sp = 0
            if (skill.rank) {
                rank = skill.rank.value
            }
            if (selectedSkillLevel == -1 && skill.skill) {
                continue
            }
            if (!skill.skill && selectedSkillLevel >= 0) {
                continue
            }
            let right = []
            if (skill.skill) {
                if (selectedSkillLevel > 0 && selectedSkillLevel !== skill.skill.trained_skill_level) {
                    continue
                }
                sp = skill.skill.skillpoints_in_skill
                for (let i = 1; i <= skill.skill.trained_skill_level; i++) {
                    right.push(`<i class="me-1 fas fa-square"></i>`)
                }
                if (skill.skill.trained_skill_level < 5) {
                    for (let i = skill.skill.trained_skill_level + 1; i <= 5; i++) {
                        right.push(`<i class="me-1 far fa-square"></i>`)
                    }
                }

            }

            const left = `
                <div>
                    ${typeName}<br>
                    <span class="text-muted">
                        Rank: ${rank} | Total SP: ${numeral(sp).format("0,0")} of ${numeral(potentialSP(rank)).format("0,0")}
                    </span>
                </div>
            `

            const out = `
            <li class="list-group-item text-white">
                <div class="d-flex w-100 justify-content-between">
                    ${left}
                    <div>
                    ${right.join("")}
                    </div>
                </div>
            </li>
            `

            skillLines.push(out)
        }

        if (skillLines.length == 0) {
            console.log("skillLines.length == 0")
            continue
        }

        cards.push(`
            <div class="card mb-2">
                <div class="card-header">
                    <h5 class="mb-0 d-flex w-100 justify-content-between">
                        <span>${group.name}</span>
                        <span>
                            (${hasSkillCount(group)} / ${numeral(group.totalGroupSP).format("0,0")})
                        </span>
                    </h5>
                </div>
                <ul class="list-group list-group-flush">
                    <li class="list-group-item text-white disabled">
                        ${levelVSkillCount(group)} Level V Skills injected for a total of ${numeral(levelVTotalSP(group)).format("0,0")}
                    </li>
                    ${skillLines.join("")}
                </ul>
            </div>
        `
        )
    }
    document.querySelector("#skillGroupDetails").innerHTML = `<ul class="list-group">${cards.join("")}</ul>`;
}
renderSkillLevelListGroup()
renderSkillGroupListGroup()