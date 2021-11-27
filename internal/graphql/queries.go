package graphql

const skillboard = `
query($id: Uint64!) {
  user(id: $id) {
    id
    character {
      id
      name
      birthday
      securityStatus
      corporation {
        id
        name
        alliance {
          id
          name
        }
      }
    }
    skillMeta {
      totalSp
      unallocatedSp
    }
    skills {
      totalGroupSP
      info {
        id
        name
      }
      skills {
        skillID
        activeSkillLevel
        trainedSkillLevel
        skillpointsInSkill
        info {
          id
          name
          attributes {
            attributeID
            value
          }
        }
      }
    }
    implants {
      implant {
        id
        name
      }
    }
    attributes {
      intelligence
      perception
      charisma
      willpower
      memory
      bonusRemaps
    }
    queue {
      skillID
      info {
        name
      }
      queuePosition
      startDate
      finishDate
    }
  }
}`
